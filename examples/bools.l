def T \a b . a
def F \a b . b

T a b
F a b

def OR \a b . a T b
def AND \a b . a b F

OR T T
OR T F
OR F T
OR F F

AND T T
AND T F
AND F T
AND F F

def IF \ b t f . b t f