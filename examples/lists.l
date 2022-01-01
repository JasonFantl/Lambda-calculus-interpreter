def Nil Pair True True
def IsNil First
def Cons \ h t . Pair False (Pair h t)
def Head \ z . First (Second z)
def Tail \z. Second (Second z)