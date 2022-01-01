def SUCC \ n f x . f (n f x)

def PLUS \ a b . a SUCC b

def MULT \ m n f . m (n f)

def POW \ m n . n m 

def PRED \ n f x . n (\ g h.h (g f)) (\u.x) (\u.u)

def SUB \m n . n PRED m

def Zero \ f x . x
def One \f x . f x
def Two SUCC One
def Three SUCC Two
def Four PLUS One Three
def Five SUCC Four

