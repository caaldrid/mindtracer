from django.urls import path
from . import views

urlpatterns = [
    path('', views.index, name='index'),
    path('new-line', views.new_line, name="newLine")
]
