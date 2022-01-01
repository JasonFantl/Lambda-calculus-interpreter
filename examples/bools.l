def True \a b . a
def False \a b . b

def NOT \ a . a False True
def OR \a b . a True b
def AND \a b . a b False

def IF \ b t f . b t f