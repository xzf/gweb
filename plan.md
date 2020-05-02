 - gweb   
      - api   
         - api call  
          in-> get func by path , ,get para  
		 call func with para   
	  - db  
		table -> struct   
	  - file  
      css,js   
 	  download file   
 	  - 伪静态   
 
		
obj -> func map 
    sample para -> call one func with para 
request path index func map 
call func 


api func para type 

func() // no para

//baseType mean go base type ,unsupport type : complex、chan、ptr and their slice type  
func(baseType)
func(baseType,baseType)
func(baseType,baseType,baseType)//Not recommended too much baseType para
...
func(baseType,...,baseType)   

//strut any level field can't be ptr 
func (struct)

!!!!
can't be func (struct,baseType)  or func (baseType,struct)
