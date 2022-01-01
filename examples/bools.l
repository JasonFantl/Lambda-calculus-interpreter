def T \a b . a
def F \a b . b

def NOT \ a . a F T
def OR \a b . a T b
def AND \a b . a b F

def IF \ b t f . b t f