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
func (struct) //strut any level field can't be ptr 

//because golang reflect can get func para name, golang base type unsupport 

context response write and request get solution: 
set respWriter and request to a global map(with lock,data race) 
key of this map is goroutine id   
and then call like code 
code in implement:
```
func(a api)Path(para Para){
   a.WriteBody(`{ "code",}`)
}
```
code in gweb package:
```
func(a WebApi)WriteBody(str string){// no must to be string para
    ctx,ok := a.globalMap[a.getGoroutineId()]
    if !ok{
        return 
    }
    ctx.writer.Write([]byte(str))
}
```
todo 
file upload 
file handler
https 