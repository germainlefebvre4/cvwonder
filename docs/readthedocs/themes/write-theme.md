# Write a theme

## Directory structure

```tree
themes
└── <my-theme-name>
    └── index.html
```

## Build the theme

The templated file is `themes/<my-theme-name>/index.html` and stands as the entry point for the theme.

The templating engine is the [`template/html` go package](https://pkg.go.dev/html/template).

The input CV data is formerly structured to make it easy to use in the template. In order to help you write the data, a JSON Schema is provided in the `schema` directory.

For example, the dv data contains the details of the CV owner:

```yaml
[...]

person:
  name: Germain
  profession: Bâtisseur de Plateformes et de Nuages

[...]
```

Which can be used in the template like this:

```html
<h1>{{ .Person.Name }}</h1>
<h2>{{ .Person.Profession }}</h2>
```

!!! note "Go template variable name"
    Every key in the yaml file is capitalized in the Go template.
    Even though the yaml file uses `person`, the Go template variable name is `Person`.
