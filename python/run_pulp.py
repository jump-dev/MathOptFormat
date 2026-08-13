# Copyright (c) 2020: Oscar Dowson and contributors
#
# Use of this source code is governed by an MIT-style license that can be found
# in the LICENSE.md file or at https://opensource.org/licenses/MIT.

import json
import pulp
import pytest
import os

ROOT = os.path.dirname(os.path.abspath(__file__)) + "/../"

class UnsupportedObjective(Exception):
    def __init__(self, F):
        self.F = F
        return

class UnsupportedConstraint(Exception):
    def __init__(self, F, S):
        self.F, self.S = F, S
        return

def parse_func(f, var_by_name):
    match f["type"]:
        case "Variable":
            return var_by_name[f["name"]]
        case "ScalarAffineFunction":
            expr = pulp.LpAffineExpression()
            for term in f.get("terms", []):
                expr += term["coefficient"] * var_by_name[term["variable"]]
            if "constant" in f:
                expr += f["constant"]
            return expr
        case _:
            return None

def add_variable_constraint(prob, f, s, var_by_name):
    x = var_by_name[f["name"]]
    match s["type"]:
        case "LessThan":
            x.upBound = s["upper"]
        case "GreaterThan":
            x.lowBound = s["lower"]
        case "EqualTo":
            x.lowBound = x.upBound = s["value"]
        case "Interval":
            x.lowBound, x.upBound = s["lower"], s["upper"]
        case "ZeroOne":
            x.cat = pulp.LpBinary
        case "Integer":
            x.cat = pulp.LpInteger
        case _:
            raise UnsupportedConstraint(f["type"], s["type"])
    return

def add_constraint(prob, c, var_by_name):
    f, s = c["function"], c["set"]
    match f["type"]:
        case "Variable":
            add_variable_constraint(prob, f, s, var_by_name)
        case _:
            expr = parse_func(f, var_by_name)
            if expr is None:
                raise UnsupportedConstraint(f["type"], s["type"])
            match s["type"]:
                case "LessThan":
                    prob += (expr <= s["upper"])
                case "GreaterThan":
                    prob += (expr >= s["lower"])
                case "EqualTo":
                    prob += (expr == s["value"])
                case _:
                    raise UnsupportedConstraint(f["type"], s["type"])
    return

def read_from_file(filename):
    with open(filename, "r") as f:
        data = json.load(f)
    return read_from_dict(data)

def read_from_dict(data):
    name = data.get("name", "")
    obj = data.get("objective", {})
    is_max = obj.get("sense", "MIN_SENSE") == "MAX_SENSE"
    sense = pulp.LpMaximize if is_max else pulp.LpMinimize
    prob = pulp.LpProblem(name=name, sense=sense)
    prob.var_by_name = {}
    for x in data.get("variables", []):
        name = x["name"]
        prob.var_by_name[name] = prob.add_variable(name, None, None)
    if "function" in obj:
        obj_f = parse_func(obj["function"], prob.var_by_name)
        if obj_f is None:
            raise UnsupportedObjective(obj["function"]["type"])
        prob += obj_f
    for c in data.get("constraints", []):
        add_constraint(prob, c, prob.var_by_name)
    return prob

def solve_from_file(filename):
    prob = read_from_file(filename)
    print(prob)
    prob.solve()
    return {
        "status": pulp.LpStatus[prob.status],
        "primal": {k: v.value() for (k, v) in prob.var_by_name.items()}
    }

# Usage:

def test_success():
    ret = solve_from_file(ROOT + 'examples/milp.mof.json')
    assert ret["status"] == "Optimal"
    assert ret["primal"] == {'x': 0.0, 'y': 1.0}
    return

def test_failure():
    with pytest.raises(UnsupportedObjective):
        solve_from_file(ROOT + 'examples/nlp.mof.json')
    return
