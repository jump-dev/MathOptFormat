import json
import jsonschema
import os


SCHEMA_FILENAME = '../schemas/mof.1.4.schema.json'

def validate(filename):
    with open(filename, 'r', encoding='utf-8') as io:
        instance = json.load(io)
    with open(SCHEMA_FILENAME, 'r', encoding='utf-8') as io:
        schema = json.load(io)
    jsonschema.validate(instance = instance, schema = schema)

def summarize_schema():
    with open(SCHEMA_FILENAME, 'r', encoding='utf-8') as io:
        schema = json.load(io)
    summary = "## Sets\n"
    summary += "\n### Scalar Sets\n\n" + summarize(schema, "scalar_sets")
    summary += "\n### Vector Sets\n\n" + summarize(schema, "vector_sets")
    summary += "\n## Functions\n"
    summary += "\n### Scalar Functions\n\n" + summarize(schema, "scalar_functions")
    summary += "\n### Vector Functions\n\n" + summarize(schema, "vector_functions")
    operators, leaves = summarize_nonlinear(schema)
    summary += "\n### Nonlinear\n"
    summary += "\n#### Leaf nodes\n\n" + leaves
    summary += '\n#### Operators\n\n' + operators
    return summary

def oneOf_to_object(item):
    head = item["properties"]["type"]
    ret = []
    if "const" in head:
        description = item.get("description", "").replace("|", "\\|")
        example = item.get("examples", [""])
        assert(len(example) == 1)
        ret.append({
            'name': head["const"],
            'description': description,
            'example': example[0],
        })
    else:
        for k in head["enum"]:
            ret.append({
                'name': k,
                'description': "",
                'example': "",
            })
    return ret

def summarize(schema, key):
    summary = \
        "| Name | Description | Example |\n" + \
        "| ---- | ----------- | ------- |\n"
    for item in schema["definitions"][key]["oneOf"]:
        for obj in oneOf_to_object(item):
            summary += "| `\"%s\"` | %s | %s |\n" % \
                (obj['name'], obj['description'], obj['example'])
    return summary

def summarize_nonlinear(schema):
    operators = "| Name | Arity |\n| ---- | ----- |\n"
    leaves = "| Name | Description | Example |\n| ---- | ----------- | ------- |\n"
    description_map = {
        "Unary operators": 'Unary',
        "Binary operators": 'Binary',
        "N-ary operators": 'N-ary',
    }
    for item in schema["definitions"]["NonlinearTerm"]["oneOf"]:
        desc = description_map.get(item["description"], None)
        if desc == None:
            obj = oneOf_to_object(item)[0]
            leaves += "| `\"%s\"` | %s | %s |\n" % \
                (obj['name'], obj['description'], obj['example'])
        else:
            for obj in oneOf_to_object(item):
                operators += "| `\"%s\"` | Unary |\n" % obj['name']
    return operators, leaves

###
### Validate all the files in the examples directory.
###

for filename in os.listdir('../examples'):
    validate(os.path.join('../examples', filename))

###
### Summarize the schema for the README table.
###

print(summarize_schema())
