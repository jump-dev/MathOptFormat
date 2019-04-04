# MathOptFormat

This repository describes a file-format for mathematical optimization problems
called _MathOptFormat_ with the file extension `.mof.json`.

It is heavily inspired by [MathOptInterface](https://github.com/JuliaOpt/MathOptInterface.jl).

## Standard form

MathOptFormat is a generic file format for mathematical optimization problems
encoded in the form

       min/max: fᵢ(x)       i=1,2,…,I
    subject to: gⱼ(x) ∈ Sⱼ  j=1,2,…,J

where `x ∈ ℝᴺ`, `fᵢ: ℝᴺ → ℝ`, and `gⱼ: ℝᴺ → ℝᴹʲ`, and `Sⱼ ⊆ ℝᴹʲ`.

The functions `fᵢ` and `gⱼ`, and sets `Sⱼ` supported by MathOptFormat are
defined in the [MathOptFormat schema](#the-schema).

The current list of supported functions and sets is not exhaustive. It is
intended that MathOptFormat will be extended in future versions to support
additional functions and sets.

## An example

The standard form described above is very general. To give a concrete example,
consider the following linear program:

           min: 2x + 1
    subject to: x ≥ 1

Encoded in our standard form, we have

    f₁(x) = 2x + 1
    g₁(x) = x
    S₁    = [1, ∞)

Encoded into the MathOptFormat file format, this example becomes:
```json
{
    "version": 1,
    "variables": [{"name": "x"}],
    "objectives": [{
        "sense": "min",
        "function": {
            "head": "ScalarAffineFunction",
            "terms": [
                {"coefficient": 2, "variable": "x"}
            ],
            "constant": 1
        }
    }],
    "constraints": [{
        "name": "x >= 1",
        "function": {"head": "SingleVariable", "variable": "x"},
        "set": {"head": "GreaterThan", "lower": 1}
    }]
}
```

Let us now describe each part of the file format in turn. First, notice that the
file format is a valid JSON (JavaScript Object Notation) file. This enables the
MathOptFormat to be both _human-readable_, and _machine-readable_. Some readers
may argue that JSON is tricky to edit as a human due to the quantity of brackets
and "boiler-plate" such as `"function"` and `"head"`. We do not disagree.
However, we believe that JSON strikes a fine balance between human and machine
readability.

Inside the document, the model is stored as a single JSON object. JSON objects
are key-value mappings enclosed by curly braces (`{` and `}`). There are four
required keys at the top level:

 - `"version"`

   An integer describing the minimum version of MathOptFormat needed to parse
   the file. This is included to safeguard against later revisions.

 - `"variables"`

   A list of JSON objects, with one object for each variable in the model. Each
   object has a required key `"name"` which maps to a unique string for that
   variable. It is illegal to have two variables with the same name. These names
   will be used later in the file to refer to each variable.

 - `"objectives"`

   A list of JSON objects, with one element for each objective in the model.
   Each object has two required keys:

    - `"sense"`

      A string which must be `"min"` or `"max"`.

    - `"function"`

      A JSON object that describes the function. There are many different types
      of functions that MathOptFormat recognizes (see
      [List of supported functions](#list-of-supported-functions)), each of
      which has a different structure. However, each function has a required key
      called `"head"` which is used to describe the type of the function. In
      this case, the function is `"ScalarAffineFunction"`.

      A `"ScalarAffineFunction"` is the function `f(x) = aᵀx + b`, where `a` is
      a constant `N×1` vector, and `b` is a scalar constant. In addition to
      `"head"`, it has two required keys:

       - `"terms"`

         A list of JSON objects, containing one object for each non-zero element
         in `a`. Each object has two required keys: `"coefficient"`, and
         `"variable"`. `"coefficient"` maps to a real number that is the
         coefficient in `a` corresponding to the variable (identified by its
         string name) in `"variable"`.

      - `"constant"`

        The value of `b`.

 - `"constraints"`

   A list of JSON objects, with one element for each constraint in the model.
   Each object has three required fields:

    - `"name"`

      A unique string name for the constraint.

    - `"function"`

      A JSON object that describes the function `gⱼ` associated with constraint
      `j`. The function field is similar to the function field in
      `"objectives"`; however, in this example, our function is a single
      variable function of the variable `"x"`.

    - `"set"`

      A JSON object that describes the set `Sⱼ` associated with constraint `j`.
      In this example, the set `[1, ∞)` is the MathOptFormat set `GreaterThan`
      with a lower bound of `1`. See [List of supported sets](#list-of-supported-sets)
      for other sets supported by MathOptFormat.

### Other examples

A number of examples of optimization problems encoded using MathOptFormat are
provided in the [`/examples` directory](https://github.com/odow/MathOptFormat/tree/master/examples).

## The schema

A [JSON schema](http://json-schema.org/) for the `.mof.json` file-format is
provided in the file [`mof.schema.json`](https://github.com/odow/MathOptFormat/blob/master/mof.schema.json).

It is intended for the schema to be self-documenting. Instead of modifying or
adding to this documentation, clarifying edits should be made to the
`description` field of the relevant part of the schema.

### List of supported functions

The list of functions supported by MathOptFormat are contained in the
`#/definitions/scalar_functions` and `#/definitions/vector_functions` fields of
the schema. Scalar functions are functions for which `Mj=1`, while vector
functions are functions for which `Mj≥1`.

Here is a summary of the functions defined by MathOptFormat.

#### Scalar Functions

| Name | Description | Example |
| ---- | ----------- | ------- |
| `"SingleVariable"` | The scalar variable `variable`. | {"head": "SingleVariable", "variable": "x"} |
| `"ScalarAffineFunction"` | The function `a'x + b`, where `a` is a sparse vector specified by a list of `ScalarAffineTerm`s in `terms` and `b` is the scalar in `constant`. Duplicate variables in `terms` are accepted, and the corresponding coefficients are summed together. | {"head": "ScalarAffineFunction", "constant": 1.0, "terms": [{"coefficient": 2.5, "variable": "x"}]} |
| `"ScalarQuadraticFunction"` | The function `0.5x'Qx + a'x + b`, where `a` is a sparse vector of `ScalarAffineTerm`s in `affine_terms`, `b` is the scalar `constant`, and `Q` is a symmetric matrix specified by a list of `ScalarQuadraticTerm`s in `quadratic_terms`. Duplicate indices in `affine_terms` and `quadratic` are accepted, and the corresponding coefficients are summed together. Mirrored indices in `quadratic_terms` (i.e., `(i,j)` and `(j, i)`) are considered duplicates; only one need to be specified. | {"head": "ScalarAffineFunction", "constant": 1.0, "affine_terms": [{"coefficient": 2.5, "variable": "x"}], "quadratic_terms": [{"coefficient": 2.0, "variable_1": "x", "variable_2": "y"}]} |
| `"Nonlinear"` | An expression graph representing a scalar function. |  |

#### Vector Functions

| Name | Description | Example |
| ---- | ----------- | ------- |
| `"VectorOfVariables"` | An ordered list of variables. | {"head": "VectorOfVariables", "variables": ["x", "y"]} |
| `"VectorAffineFunction"` | The function `Ax + b`, where `A` is a sparse matrix specified by a list of `VectorAffineTerm`s in `terms` and `b` is a dense vector specified by `constants`. | {"head": "VectorAffineFunction", "constants": [1.0], "terms": [{"output_index": 1, "scalar_term": {"coefficient": 2.5, "variable": "x"}}]} |
| `"VectorQuadraticFunction"` | The vector-valued quadratic function `q(x) + Ax + b`, where `q(x)` is specified by a list of `VectorQuadraticTerm`s in `quadratic_terms`, `A` is a sparse matrix specified by a list of `VectorAffineTerm`s in `affine_terms` and `b` is a dense vector specified by `constants`. |  |

### List of supported sets

The list of sets supported by MathOptFormat are contained in the
`#/definitions/scalar_sets` and `#/definitions/vector_sets` fields of the
schema. Scalar sets are sets for which `Mj=1`, while vector sets are sets for
which `Mj≥1`.

Here is a summary of the sets defined by MathOptFormat.

#### Scalar Sets

| Name | Description | Example |
| ---- | ----------- | ------- |
| `"LessThan"` | (-∞, upper] | {"head": "LessThan", "upper": 2.1} |
| `"GreaterThan"` | [lower, ∞) | {"head": "GreaterThan", "lower": 2.1} |
| `"EqualTo"` | {value} | {"head": "EqualTo", "value": 2.1} |
| `"Interval"` | [lower, upper] | {"head": "Interval", "lower": 2.1, "upper": 3.4} |
| `"Semiinteger"` | {0} ∪ {lower, lower + 1, ..., upper} | {"head": "Semiinteger", "lower": 2, "upper": 4} |
| `"Semicontinuous"` | {0} ∪ [lower, upper] | {"head": "Semicontinuous", "lower": 2.1, "upper": 3.4} |
| `"ZeroOne"` | {0, 1} | {"head": "ZeroOne"} |
| `"Integer"` | ℤ | {"head": "Integer"} |

#### Vector Sets

| Name | Description | Example |
| ---- | ----------- | ------- |
| `"ExponentialCone"` | [x, y, z] ∈ {R³: y * exp(x / y) ≤ z, y ≥ 0} | {"head": "ExponentialCone"} |
| `"DualExponentialCone"` | [u, v, w] ∈ {R³: -u * exp(v / u) ≤ exp(1) * w, u < 0} | {"head": "DualExponentialCone"} |
| `"SOS1"` | A special ordered set of type I. | {"head": "SOS1", "weights": [1, 3, 2]} |
| `"SOS2"` | A special ordered set of type II. | {"head": "SOS2", "weights": [1, 3, 2]} |
| `"GeometricMeanCone"` | [t, x] ∈ {R^{dimension}: t ≤ (Πxᵢ)^{1 / (dimension-1)}} | {"head": "GeometricMeanCone", "dimension": 3} |
| `"SecondOrderCone"` | [t, x] ∈ {R^{dimension} : t ≥ \|\|x\|\|₂ | {"head": "SecondOrderCone", "dimension": 3} |
| `"RotatedSecondOrderCone"` | [t, u, x] ∈ {R^{dimension} : 2tu ≥ (\|\|x\|\|₂)²; t, u ≥ 0} | {"head": "RotatedSecondOrderCone", "dimension": 3} |
| `"Zeros"` | {0}^{dimension} | {"head": "Zeros", "dimension": 3} |
| `"Reals"` | R^{dimension} | {"head": "Reals", "dimension": 3} |
| `"Nonpositives"` | R₋^{dimension} | {"head": "Nonpositives", "dimension": 3} |
| `"Nonnegatives"` | R₊^{dimension} | {"head": "Nonnegatives", "dimension": 3} |
| `"RootDetConeTriangle"` |  |  |
| `"RootDetConeSquare"` |  |  |
| `"LogDetConeTriangle"` |  |  |
| `"LogDetConeSquare"` |  |  |
| `"PositiveSemidefiniteConeTriangle"` |  |  |
| `"PositiveSemidefiniteConeSquare"` |  |  |
| `"PowerCone"` | [x, y, z] ∈ {R³: x^{exponent} y^{1-exponent} ≥ \|z\|; x, y ≥ 0} | {"head": "PowerCone", "exponent": 2.0} |
| `"DualPowerCone"` | [u, v, w] ∈ {R³: (u / exponent)^{exponent} (v / (1-exponent))^{1-exponent} ≥ \|w\|; u, v ≥ 0} | {"head": "DualPowerCone", "exponent": 2.0} |

## Implementations

- Julia

  - The [MathOptFormat.jl](https://github.com/odow/MathOptFormat.jl) package
    supports reading and writing MathOptFormat files.
