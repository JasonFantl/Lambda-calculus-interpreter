def SUCC \ n f x . f (n f x)

def Zero \ f x . x
def One \f x . f x
def Two SUCC One
def Three SUCC Two

def PLUS \ a b . a SUCC b