def IsZero \n . n (\x. False) True
def LEQ \m n . IsZero (SUB m n)
def EQ \ m n . AND (LEQ m n) (LEQ n m)