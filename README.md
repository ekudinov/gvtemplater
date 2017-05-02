# gvtemplater
Tool for generation go file from vue templates.

  This program is designed to generate a go file from vue templates.
The reason for creating this program is to write vue templates with native ide support,
and then create from them go file in which the templates are represented as constants
and easily included in the program code developed under the library gopherjs-vue
https://github.com/oskca/gopherjs-vue.
  For example, a template named button.vue turns into a constant button = "" etc.
And then in the right place this constant can be used for its intended purpose.

Project structure for example before generation:

![Alt text](pics/start.png)

Vue template in IDE:

![Alt text](pics/template.png)

Generated go file:

![Alt text](pics/generated.png)

Project structure after generation:

![Alt text](pics/end.png)

At this point, the program checks for duplicate template names, empty templates,
and if there are no templates with a .vue extension.

Help options:

![Alt text](pics/help.png)
