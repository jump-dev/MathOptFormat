# MathOptFormat

This repository describes a file-format for mathematical optimization problems
called _MathOptFormat_ with the file extension `.mof.json`.

MathOptFormat is rigidly defined by the [JSON schema](http://json-schema.org/)
available at
[`https://jump.dev/MathOptFormat/mof.0.5.schema.json`](https://jump.dev/MathOptFormat/schemas/mof.0.5.schema.json).

It is intended for the schema to be self-documenting. Instead of modifying or
adding to this documentation, clarifying edits should be made to the
`description` field of the relevant part of the schema.

A number of examples of optimization problems encoded using MathOptFormat are
provided in the [`/examples` directory](https://github.com/odow/MathOptFormat/tree/master/examples).

A paper describing the motivation, design principles, and historical setting of
MathOptFormat is available at:

Legat, B., Dowson, O., Garcia, J.D., Lubin, M. (2020). MathOptInterface: a data
structure for mathematical optimization problems.
[[preprint]](http://www.optimization-online.org/DB_HTML/2020/02/7609.html)
[[repository]](https://github.com/jump-dev/MathOptFormat)

**We highly recommend you read that paper before reading further.**

## Implementations

- Julia

  - The [MathOptInterface.jl](https://github.com/jump-dev/MathOptInterface.jl) package
    supports reading and writing MathOptFormat files.

## Standard form

MathOptFormat is a generic file format for mathematical optimization problems
encoded in the form

       min/max: f₀(x)
    subject to: fᵢ(x) ∈ Sᵢ  i=1,2,…,I

where `x ∈ ℝᴺ`, `fᵢ: ℝᴺ → ℝᴹⁱ`, and `Sᵢ ⊆ ℝᴹⁱ`.

The functions `fᵢ` and sets `Sᵢ` supported by MathOptFormat are defined in the
[MathOptFormat schema](#the-schema).

The current list of supported functions and sets is not exhaustive. It is
intended that MathOptFormat will be extended in future versions to support
additional functions and sets.

## An example

The standard form described above is very general. To give a concrete example,
consider the following linear program:

           min: 2x + 1
    subject to: x ≥ 1

Encoded in our standard form, we have

    f₀(x) = 2x + 1
    f₁(x) = x
    S₁    = [1, ∞)

Encoded into the MathOptFormat file format, this example becomes:
```json
{
    "version": {
        "major": 0,
        "minor": 5
    },
    "variables": [{"name": "x"}],
    "objective": {
        "sense": "min",
        "function": {
            "type": "ScalarAffineFunction",
            "terms": [
                {"coefficient": 2, "variable": "x"}
            ],
            "constant": 1
        }
    },
    "constraints": [{
        "name": "x >= 1",
        "function": {"type": "SingleVariable", "variable": "x"},
        "set": {"type": "GreaterThan", "lower": 1}
    }]
}
```

Let us now describe each part of the file format in turn. First, notice that the
file format is a valid JSON (JavaScript Object Notation) file. This enables the
MathOptFormat to be both _human-readable_, and _machine-readable_. Some readers
may argue that JSON is tricky to edit as a human due to the quantity of brackets
and "boiler-plate" such as `"function"` and `"type"`. We do not disagree.
However, we believe that JSON strikes a fine balance between human and machine
readability.

Inside the document, the model is stored as a single JSON object. JSON objects
are key-value mappings enclosed by curly braces (`{` and `}`). There are four
required keys at the top level:

 - `"version"`

   An object describing the minimum version of MathOptFormat needed to parse
   the file. This is included to safeguard against later revisions. It contains
   two fields: `"major"` and `"minor"`. These fields should be interpreted
   using [SemVer](https://semver.org).

 - `"variables"`

   A list of JSON objects, with one object for each variable in the model. Each
   object has a required key `"name"` which maps to a unique string for that
   variable. It is illegal to have two variables with the same name. These names
   will be used later in the file to refer to each variable.

 - `"objective"`

   A JSON objects describing the objective of the model. It has one required
   keys:

    - `"sense"`

      A string which must be `"min"`, `"max"`, or `"feasibility"`. If the sense
      is `min` or `max`, a second key `"function"`, must be defined:

    - `"function"`

      A JSON object that describes the function. There are many different types
      of functions that MathOptFormat recognizes (see
      [List of supported functions](#list-of-supported-functions)), each of
      which has a different structure. However, each function has a required key
      called `"type"` which is used to describe the type of the function. In
      this case, the function is `"ScalarAffineFunction"`.

      A `"ScalarAffineFunction"` is the function `f(x) = aᵀx + b`, where `a` is
      a constant `N×1` vector, and `b` is a scalar constant. In addition to
      `"type"`, it has two required keys:

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

### List of supported functions

The list of functions supported by MathOptFormat are contained in the
`#/definitions/scalar_functions` and `#/definitions/vector_functions` fields of
the schema. Scalar functions are functions for which `Mi=1`, while vector
functions are functions for which `Mi≥1`.

Here is a summary of the functions defined by MathOptFormat.

#### Scalar Functions

| Name | Description | Example |
| ---- | ----------- | ------- |
| `"SingleVariable"` | The scalar variable `variable`. | {"type": "SingleVariable", "variable": "x"} |
| `"ScalarAffineFunction"` | The function `a'x + b`, where `a` is a sparse vector specified by a list of `ScalarAffineTerm`s in `terms` and `b` is the scalar in `constant`. Duplicate variables in `terms` are accepted, and the corresponding coefficients are summed together. | {"type": "ScalarAffineFunction", "constant": 1.0, "terms": [{"coefficient": 2.5, "variable": "x"}]} |
| `"ScalarQuadraticFunction"` | The function `0.5x'Qx + a'x + b`, where `a` is a sparse vector of `ScalarAffineTerm`s in `affine_terms`, `b` is the scalar `constant`, and `Q` is a symmetric matrix specified by a list of `ScalarQuadraticTerm`s in `quadratic_terms`. Duplicate indices in `affine_terms` and `quadratic` are accepted, and the corresponding coefficients are summed together. Mirrored indices in `quadratic_terms` (i.e., `(i,j)` and `(j, i)`) are considered duplicates; only one need to be specified. | {"type": "ScalarQuadraticFunction", "constant": 1.0, "affine_terms": [{"coefficient": 2.5, "variable": "x"}], "quadratic_terms": [{"coefficient": 2.0, "variable_1": "x", "variable_2": "y"}]} |
| `"ScalarNonlinearFunction"` | An expression graph representing a scalar nonlinear function. |  |

For more information on `"ScalarNonlinearFunction"` functions, see
[Nonlinear functions](nonlinear-functions).

#### Vector Functions

| Name | Description | Example |
| ---- | ----------- | ------- |
| `"VectorOfVariables"` | An ordered list of variables. | {"type": "VectorOfVariables", "variables": ["x", "y"]} |
| `"VectorAffineFunction"` | The function `Ax + b`, where `A` is a sparse matrix specified by a list of `VectorAffineTerm`s in `terms` and `b` is a dense vector specified by `constants`. | {"type": "VectorAffineFunction", "constants": [1.0], "terms": [{"output_index": 1, "scalar_term": {"coefficient": 2.5, "variable": "x"}}]} |
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
| `"LessThan"` | (-∞, upper] | {"type": "LessThan", "upper": 2.1} |
| `"GreaterThan"` | [lower, ∞) | {"type": "GreaterThan", "lower": 2.1} |
| `"EqualTo"` | {value} | {"type": "EqualTo", "value": 2.1} |
| `"Interval"` | [lower, upper] | {"type": "Interval", "lower": 2.1, "upper": 3.4} |
| `"Semiinteger"` | {0} ∪ {lower, lower + 1, ..., upper} | {"type": "Semiinteger", "lower": 2, "upper": 4} |
| `"Semicontinuous"` | {0} ∪ [lower, upper] | {"type": "Semicontinuous", "lower": 2.1, "upper": 3.4} |
| `"ZeroOne"` | {0, 1} | {"type": "ZeroOne"} |
| `"Integer"` | ℤ | {"type": "Integer"} |

#### Vector Sets

| Name | Description | Example |
| ---- | ----------- | ------- |
| `"ExponentialCone"` | [x, y, z] ∈ {R³: y * exp(x / y) ≤ z, y ≥ 0} | {"type": "ExponentialCone"} |
| `"DualExponentialCone"` | [u, v, w] ∈ {R³: -u * exp(v / u) ≤ exp(1) * w, u < 0} | {"type": "DualExponentialCone"} |
| `"SOS1"` | A special ordered set of type I. | {"type": "SOS1", "weights": [1, 3, 2]} |
| `"SOS2"` | A special ordered set of type II. | {"type": "SOS2", "weights": [1, 3, 2]} |
| `"GeometricMeanCone"` | [t, x] ∈ {R^{dimension}: t ≤ (Πxᵢ)^{1 / (dimension-1)}} | {"type": "GeometricMeanCone", "dimension": 3} |
| `"SecondOrderCone"` | [t, x] ∈ {R^{dimension} : t ≥ \|\|x\|\|₂ | {"type": "SecondOrderCone", "dimension": 3} |
| `"RotatedSecondOrderCone"` | [t, u, x] ∈ {R^{dimension} : 2tu ≥ (\|\|x\|\|₂)²; t, u ≥ 0} | {"type": "RotatedSecondOrderCone", "dimension": 3} |
| `"Zeros"` | {0}^{dimension} | {"type": "Zeros", "dimension": 3} |
| `"Reals"` | R^{dimension} | {"type": "Reals", "dimension": 3} |
| `"Nonpositives"` | R₋^{dimension} | {"type": "Nonpositives", "dimension": 3} |
| `"Nonnegatives"` | R₊^{dimension} | {"type": "Nonnegatives", "dimension": 3} |
| `"RootDetConeTriangle"` | {[t, X] ∈ R^{1 + d(d+1)/2} : t ≤ det(X)^{1/d}}, where the matrix `X` is represented in the same symmetric packed format as in the `PositiveSemidefiniteConeTriangle`. The argument `side_dimension` is the side dimension of the matrix `X`, i.e., its number of rows or columns. | {"type": "RootDetConeTriangle", "side_dimension": 2} |
| `"RootDetConeSquare"` | {[t, X] ∈ R^{1 + d^2} : t ≤ det(X)^{1/d}, X symmetric}, where the matrix `X` is represented in the same symmetric packed format as in the `PositiveSemidefiniteConeSquare`. The argument `side_dimension` is the side dimension of the matrix `X`, i.e., its number of rows or columns. | {"type": "RootDetConeSquare", "side_dimension": 2} |
| `"LogDetConeTriangle"` | {[t, u, X] ∈ R^{2 + d(d+1)/2} : t ≤ u log(det(X/u)), u > 0}, where the matrix `X` is represented in the same symmetric packed format as in the `PositiveSemidefiniteConeTriangle`. The argument `side_dimension` is the side dimension of the matrix `X`, i.e., its number of rows or columns. | {"type": "LogDetConeTriangle", "side_dimension": 2} |
| `"LogDetConeSquare"` | {[t, u, X] ∈ R^{2 + d^2} : t ≤ u log(det(X/u)), X symmetric, u > 0}, where the matrix `X` is represented in the same symmetric packed format as in the `PositiveSemidefiniteConeSquare`. The argument `side_dimension` is the side dimension of the matrix `X`, i.e., its number of rows or columns. | {"type": "LogDetConeSquare", "side_dimension": 2} |
| `"PositiveSemidefiniteConeTriangle"` | The (vectorized) cone of symmetric positive semidefinite matrices, with `side_dimension` rows and columns. The entries of the upper-right triangular part of the matrix are given column by column (or equivalently, the entries of the lower-left triangular part are given row by row). | {"type": "PositiveSemidefiniteConeTriangle", "side_dimension": 2} |
| `"PositiveSemidefiniteConeSquare"` | The cone of symmetric positive semidefinite matrices, with side length `side_dimension`. The entries of the matrix are given column by column (or equivalently, row by row). The matrix is both constrained to be symmetric and to be positive semidefinite. That is, if the functions in entries `(i, j)` and `(j, i)` are different, then a constraint will be added to make sure that the entries are equal. | {"type": "PositiveSemidefiniteConeSquare", "side_dimension": 2} |
| `"PowerCone"` | [x, y, z] ∈ {R³: x^{exponent} y^{1-exponent} ≥ \|z\|; x, y ≥ 0} | {"type": "PowerCone", "exponent": 2.0} |
| `"DualPowerCone"` | [u, v, w] ∈ {R³: (u / exponent)^{exponent} (v / (1-exponent))^{1-exponent} ≥ \|w\|; u, v ≥ 0} | {"type": "DualPowerCone", "exponent": 2.0} |
| `"IndicatorSet"` | If `activate_on=one`: (y, x) ∈ {0,1}×Rᴺ: y = 0 ⟹ x ∈ S, otherwise when `activate_on=zero`: (y, x) ∈ {0,1}×Rᴺ: y = 1 ⟹ x ∈ S. | {"type": "IndicatorSet", "set": {"type": "LessThan", "upper": 2.0}, "activate_on": "one"} |
| `"NormOneCone"` | (t, x) ∈ {R^{dimension}: t ≥ Σᵢ\|xᵢ\|} | {"type": "NormOneCone", "dimension": 2} |
| `"NormInfinityCone"` | (t, x) ∈ {R^{dimension}: t ≥ maxᵢ\|xᵢ\|} | {"type": "NormInfinityCone", "dimension": 2} |


### Nonlinear functions

Nonlinear functions are encoded in MathOptFormat by an expression graph. Each
expression graphs is stored in Polish prefix notation. For example, the
nonlinear expression `sin²(x)` is expressed as `^(sin(x), 2)`.

The expression graph is stored as an object with three required fields:
`"type"`, which must be `"ScalarNonlinearFunction"`, as well as `"root"` and
`"node_list"`.

`"root"` contains an object defining the root node of the expression graph. All
other nodes are stored as a flattened list in the `"node_list"` field. We
elaborate on permissible nodes and how to store them in the following
subsections.

#### Leaf nodes

Leaf nodes in the expression graph are data: they can either reference
optimization variables, or be real or complex valued numeric constants. They are
described as follows.

| Head | Description | Example |
| ---- | ----------- | ------- |
| `"real"` | A real-valued numeric constant | {"type": "real", "value": 1.0} |
| `"complex"` | A complex-valued numeric constant | {"type": "complex", "real": 1.0, "imag": 2.0} |
| `"variable"` | A reference to an optimization variable | {"type": "variable", "name": "x"} |

Nodes in the flattened list `"node_list"` can be referenced by an object with
the `"type"` field `"node"` and a field `"index"` that is the one-based index of
the node in `"node_list"`.

| Head | Description | Example |
| ---- | ----------- | ------- |
| `"node"` | A pointer to a (1-indexed) element in the `node_list` field in a nonlinear function | {"type": "node", "index": 2} |

#### Operators

All nonlinear operators in MathOptFormat are described by a JSON object with two fields:

 - `"type"`

   A string that corresponds to the operator.

 - `"args"`

   An ordered list of nodes that are passed as arguments to the operator.

The number of elements in `"args"` depends on the arity of the operator. MathOptFormat distinguishes between three arities:

 - Unary operators take one argument
 - Binary operators take two arguments
 - N-ary operators take at least one argument

To give some examples, the unary function `log(x)` is encoded as:
```json
{
  "type": "log",
  "args": [
    {"type": "variable", "name": "x"}
  ]
}
```
The binary function `x^2` (i.e., `^(x, 2)`) is encoded as:
```json
{
  "type": "^",
  "args": [
    {"type": "variable", "name": "x"},
    {"type": "real", "value": 2}
  ]
}
```
The n-ary function `x + y + 1` (i.e., `+(x, y, 1)`) is encoded as:
```json
{
  "type": "+",
  "args": [
    {"type": "variable", "name": "x"},
    {"type": "variable", "name": "y"},
    {"type": "real", "value": 1}
  ]
}
```

Here is a complete list of the nonlinear operators supported by MathOptFormat
and their corresponding arity.

| Name | Arity |
| ---- | ----- |
| `"log"` | Unary |
| `"log10"` | Unary |
| `"exp"` | Unary |
| `"sqrt"` | Unary |
| `"floor"` | Unary |
| `"ceil"` | Unary |
| `"abs"` | Unary |
| `"cos"` | Unary |
| `"sin"` | Unary |
| `"tan"` | Unary |
| `"acos"` | Unary |
| `"asin"` | Unary |
| `"atan"` | Unary |
| `"cosh"` | Unary |
| `"sinh"` | Unary |
| `"tanh"` | Unary |
| `"acosh"` | Unary |
| `"asinh"` | Unary |
| `"atanh"` | Unary |
| `"/"` | Binary |
| `"^"` | Binary |
| `"+"` | N-ary |
| `"-"` | N-ary |
| `"*"` | N-ary |
| `"min"` | N-ary |
| `"max"` | N-ary |

#### Example

As an example, consider the function `f(x, y) = (1 + 3i) * x + sin^2(x) + y`.

In Polish notation, the expression graph is:
`f(x, y) = +(*(1 + 3i, x), ^(sin(x), 2), y)`.

In MathOptFormat, this expression graph can be encoded as follows:
```json
{
  "type": "ScalarNonlinearFunction",
  "root": {
    "type": "+",
    "args": [
      {"type": "node", "index": 1},
      {"type": "node", "index": 3},
      {"type": "variable", "name": "y"}
    ]
  },
  "node_list": [
    {
      "type": "*",
      "args": [
        {"type": "complex", "real": 1, "imag": 3},
        {"type": "variable", "name": "x"}
      ]
    }, {
      "type": "sin",
      "args": [
        {"type": "variable", "name": "x"}
      ]
    }, {
      "type": "^",
      "args": [
        {"type": "node", "index": 2},
        {"type": "real", "value": 2}
      ]
    }
  ]
}
```
