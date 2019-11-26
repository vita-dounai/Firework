# Firework

## 一、 Firework LL(1)文法（未提取公因子、未消除左递归）

* 使用[...]表示终止Token
* 加粗字体表示非终止Token
* 使用π表示空串

**Program** => **Statements**

------

**Statements** => **Statement** **Statements** |  
π

------

**Statement** => **ReturnStatement** |  
**AssignStatement** |  
**WhileStatement** |  
**BlockStatement** |  
**BreakStatement** |  
**ContinueStatement** |  
**ExpressionStatement**

------

**OptionalSemicolon** => [;] |  
π

------

**ReturnStatement** => [return] **Expression** **OptionalSemicolon**

------

**AssignStatement** => [identifier] [=] **Expression** **OptionalSemicolon**

------

**WhileStatement** => [while] **Expression** **BlockStatement**

------

**BlockStatement** => [{] **Statements** [}]

------

**BreakStatement** => [break] **OptionalSemicolon**

------

**ContinueStatement** => [continue] **OptionalSemicolon**

------

**ExpressionStatement** => **Expression** **OptionalSemicolon**

------

**Expression** => [identifier] |  
[int] |  
[true] |  
[false] |  
[string] |  
**PrefixOp** **Expression** |  
**Expression** **InfixOp** **Expression** |  
**IfExpression** |  
**GroupExpression** |  
**Function** |  
**Array** |  
**Map** |  
**CallExpression** |  
**IndexExpression**

------

**PrefixOp** => [!] |  
[-]

------

**InfixOp** => [+] |  
[-] |  
[*] |  
[/] |  
[=] |  
[!=] |  
[<] |  
[>] |  
[<=] |  
[>=] |  
[**] |  
[%]

------

**IfExpression** => [if] **Expression** **BlockStatement** **Alternative**

------

**Alternative** => [else] **BlockStatement** | π

------

**GroupExpression** => [(] **Expression** [)]

------

**Function** => [|] **ParameterList** [|] **BlockStatement**

------

**ParameterList** => [identifier] **Parameters** | π

------

**Parameters** => [,] [identifier] |  
π

------

**Array** => [[] **ExpressionList** []]

------

**ExpressionList** => **Expression** **Expressions** |  
π

------

**Expressions** => [,] **Expression** |  
π

------

**Map** => [{] **PairList** [}]

------

**PairList** => **Expression**  [:] **Expression**  **Pairs** |  
π

------

**Pairs** => [,] **Expression** [:] **Expression** |  
π

------

**CallExpression** => **FunctionRef** [(] **ExpressionList** [)]

------

**FunctionRef** => **Function** |  
[identifier]

------

**IndexExpression** => **IndexableRef** [[] **Expression** []]

------

**IndexableRef** => **Map** |  
**Array** |  
[identifier]
