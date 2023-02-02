# visca
forks from ```github.com/byuoitav/visca```

Thanks to byuoitav for providing ideas and most of the logic


Added focus and discovery devices to the original.

Example of use:

```golang
func findDevices(){
   devs := Discover()
   dev := devs[0].New()
   dev.FocusFar(context.TODO(),0x00)
   ...
}
```

will add on future:

```shell
   Command: exposure,color,detail,knee...
   Inquiry Command: exposure,color,detail,knee...
   ....
```
