# zgolib
Zecchan's golang library

## Reflection package
### Map(from interface{}, to interface{}) error
Maps all members of from to matching member of to. Returns error if fail.
> Example: reflection.Map(&objFrom, &objTo)

### MapSlice(from interface{}, to interface{}) error
Maps all element of from to slice to. Returns error if fail.
> Example: reflection.MapSlice(&objFrom, &objTo)

### GetType(obj interface{}) reflect.Value, reflect.Type, bool
Gets type and value of an object. Returned bool value indicates validity of the specified obj.
> Example: reflection.GetType(&obj)

## Strformat package
### type StringFormatter
#### StringFormatter.CustomFormat  map[string]func(string) string
Create a custom formatter:
````go
 sf.CustomFormat["%hello%"] = func(txt string) string { 
		return strings.Replace(txt, "%hello%", "Hello", 1) 
 }
````
#### StringFormatter.UseCustomTime bool
Set this to true if you want to specify a custom time rather than using time.Now()
#### StringFormatter.CustomTime    time.Time
The time format that will be used for %date()% if UseCustomTime is set to true.
