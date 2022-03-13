import json
import os
from io import StringIO

import jsonpickle
from flask import Flask, redirect, url_for, flash, request
from flask_admin.babel import gettext
from flask_admin.contrib.sqla import ModelView
from flask_admin.form import FormOpts
from flask_admin.helpers import get_redirect_target
from flask_sqlalchemy import SQLAlchemy

from flask_admin import Admin, expose

import pika

app = Flask(__name__)
app.debug = True

app.config['FLASK_ENV'] = os.getenv("FLASK_ENV")
# Scheme: "postgres+psycopg2://<USERNAME>:<PASSWORD>@<IP_ADDRESS>:<PORT>/<DATABASE_NAME>"
app.config[
    'SQLALCHEMY_DATABASE_URI'] = f'postgresql://{os.getenv("POSTGRES_USER")}:{os.getenv("POSTGRES_PASSWORD")}@pg_db:5432/{os.getenv("POSTGRES_DB")}'
# app.config['SQLALCHEMY_DATABASE_URI'] = f'postgresql://flask:flask@localhost:5432/flask'
app.config['SECRET_KEY'] = 'anykey'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = True

db = SQLAlchemy(app)

connection = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq'))
channel = connection.channel()


class Category(db.Model):
    __tablename__ = "category"

    id = db.Column(db.Integer, autoincrement=True, primary_key=True)
    name = db.Column(db.String(20), unique=False, nullable=False)
    items_in_category = db.relationship('Items', uselist=False, back_populates='category_of_items', lazy=True,
                                        cascade="all")
    parent_id = db.Column(db.Integer, db.ForeignKey('category.id'), index=True)
    parent = db.relationship(lambda: Category, remote_side=id, backref='sub_category')

    def toJSON(self):
        return json.dumps(self, default=lambda o: o.toDict(),
                          sort_keys=True, indent=4)

    def __repr__(self):
        return self.name

    def toDict(self):
        return {'id': self.id, 'name': self.name, 'parent': self.parent_id}


class CategoryView(ModelView):
    column_list = (
        'name', 'parent',
    )
    form_columns = (
        'name', 'parent',
    )

    @expose('/new/', methods=('GET', 'POST'))
    def create_view(self):
        """
            Create model view
        """
        return_url = get_redirect_target() or self.get_url('.index_view')

        if not self.can_create:
            return redirect(return_url)

        form = self.create_form()
        if not hasattr(form, '_validated_ruleset') or not form._validated_ruleset:
            self._validate_form_instance(ruleset=self._form_create_rules, form=form)

        if self.validate_form(form):
            # in versions 1.1.0 and before, this returns a boolean
            # in later versions, this is the model itself
            model = self.create_model(form)

            if model:
                # TODO ivent to rabbitMQ
                channel.exchange_declare(exchange='Shopper',
                                         exchange_type='direct',auto_delete=True)
                io = r'["id":{model.id},"Name":{model.name}]'
                app.logger.info(jsonpickle.dumps(model.toDict()))
                channel.basic_publish(exchange='Shopper',
                                      routing_key='addCategory',
                                      body=jsonpickle.dumps(model.toDict()))


                # model.id model.name changed
                flash(gettext('Record was successfully created.'), 'success')
                if '_add_another' in request.form:
                    return redirect(request.url)
                elif '_continue_editing' in request.form:
                    # if we have a valid model, try to go to the edit view
                    if model is not True:
                        url = self.get_url('.edit_view', id=self.get_pk_value(model), url=return_url)
                    else:
                        url = return_url
                    return redirect(url)
                else:
                    # save button

                    return redirect(self.get_save_return_url(model, is_created=True))

        form_opts = FormOpts(widget_args=self.form_widget_args,
                             form_rules=self._form_create_rules)

        if self.create_modal and request.args.get('modal'):
            template = self.create_modal_template
        else:
            template = self.create_template

        return self.render(template,
                           form=form,
                           form_opts=form_opts,
                           return_url=return_url)


class ItemView(ModelView):
    column_list = (
        'name', "category_of_items", "price", "image_url",
    )

    form_columns = (
        'name', "category_of_items", "price", "image_url",
    )

    column_sortable_list = (
        'name', ("name", "price"), "image_url",
    )


class Items(db.Model):
    __tablename__ = "items"

    id = db.Column(db.Integer, autoincrement=True, primary_key=True)
    name = db.Column(db.String(20), unique=False, nullable=False)
    price = db.Column(db.String(120), unique=False, nullable=False)
    image_url = db.Column(db.String(200), nullable=False, default='default.jpg')
    item_category = db.Column(db.Integer, db.ForeignKey("category.id"), nullable=True)
    category_of_items = db.relationship('Category', back_populates='items_in_category', lazy=True, cascade="all")


admin = Admin(app, name="Каталог", template_mode="bootstrap3")
admin.add_view(CategoryView(Category, db.session, name="Категории"))
admin.add_view(ModelView(Items, db.session, name="Товары"))

db.create_all()


@app.route('/')
def hello_world():
    return redirect(url_for('admin.index'))
