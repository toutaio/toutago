# Toutā

## Objective
The project attemps to be a framework with great decoupling, extensibility and
maintainbility at is core.
It should be easy to deploy, easy to shape (multiple recipes and custom recipes)
and very easy to extend.
It will encourage a culture of separation of concepts, of small nemetons and
single responsability.
Since each component should be as much as possible independent of anything else, the system could allow 2 ways of implementing relationships:
- Configuration based flows: in a series of configuration files, we could define the routing of messages around the system.
- Code based flows: optionally, when the developer considers it, it could simplify the implementation by declaring the relationship in components that call other components withing the code.
Having the two options would enable the developer to decide when it worth additional decoupling or when is acceptable to declare directly the relationship.

## Nemeton concept
In order to encourage separate components, delivery and concepts, the system would allow and encourage to instead of having all the code in one place, to allow the creation of "nemetons" or pack of components, that includes everything that the developer decides to pack toguether for a small feature, or use case.
This could includes templates, engines, routes, migrations, etc.
These nemetons could be develop in a directory of the project or could be fully packed and imported in the project, in that case, having the project pull it with a nemeton manager, and locate it in in a non code directory for libraries and imported nemetons.
This should also allow to move any nemeton to and from the code directories and the imported directories, to encorage that any initially developed pack could be moved to an external repository and imported.

## Recipes concept
We want after all complete solutions. Therefore, the sytem will encourage a way to pack a complete solution either to use this as the custom project that needs to be deployed, or to offer specific assemblies for specific products. For example, a recipe for a Wiki, or a Blog, or an eCommerce, or Uncle Bob's bakery site.

## Ideas of components
- Frontend engine & post backs: using templates with some dialect and backend code, it could generate the front end output including JavaScript, DOM manipulation and reactivity to inmediately post back to a generic endpoint that could identify which page and which component is been used by the user and post back validations, fetching or changes that inmediately receive feedback from the backend. But, allowing a very simple usage of components in the templates that the developer won't need to know the actual implementation of the inners of it, but program conceptually. Example, a text field added in the template, with an postback event of "changed" that execute code in the backend.

