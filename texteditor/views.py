from django.shortcuts import render
from django.template import Template, RequestContext
from django.http import HttpResponse

# Create your views here.
def index(request):
    return render(request, 'index.html')

def new_line(request):
    content = Template(
"""<span class="line selected" hx-on:click="selectLine(event)">
    <span class="number">1</span>
    <span class="text"><span class="caret" id="caret"></span></span>
</span>"""
    )
    return HttpResponse(content.render(RequestContext(request)))