package ice

import (
	"fmt"
	"reflect"
	"strings"

	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/web"
	kit "github.com/shylinux/toolkits"
	log "github.com/shylinux/toolkits/logs"
)

func ref(obj interface{}) (reflect.Type, reflect.Value) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	return t, v
}
func val(m *ice.Message, arg ...string) []reflect.Value {
	args := []reflect.Value{reflect.ValueOf(&Message{m})}
	for _, v := range arg {
		args = append(args, reflect.ValueOf(v))
	}
	return args
}
func cmd(command *ice.Command, obj interface{}) {
	t, v := ref(obj)
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)

		if key := strings.ToLower(t.Method(i).Name); key == "list" {
			command.Hand = func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
				method.Call(val(m, arg...))
			}
		} else {
			command.Action[key] = &ice.Action{Hand: func(m *ice.Message, arg ...string) {
				method.Call(val(m, arg...))
			}}
		}
	}
}
func Cmd(key string, obj interface{}, shows ...[]*Show) {
	command := &ice.Command{Action: map[string]*ice.Action{}}
	config := &ice.Config{Value: kit.Data()}
	meta := kit.Value(config.Value, kit.MDB_META)

	show := []*Show{}
	for _, s := range shows {
		show = append(show, s...)
	}

	t, v := ref(obj)
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		if shower, ok := v.Field(i).Interface().(Shower); ok {
			show = shower.Show(key, show)
			cmd(command, v.Field(i).Interface())
			for _, k := range []string{"short", "field"} {
				kit.Value(meta, k, t.Field(i).Tag.Get(k))
			}
			continue
		}

		switch t.Field(i).Name {
		case "Name":
			command.Name = v.Field(i).String()
		case "Help":
			command.Help = v.Field(i).String()
		default:
			kit.Value(meta, strings.ToLower(t.Field(i).Name), v.Field(i).String())
		}
	}

	if shower, ok := obj.(Shower); ok {
		if kit.Format(kit.Value(meta, kit.MDB_SHORT)) == "" {
			kit.Value(meta, kit.MDB_SHORT, shower.ShortDef())
		}
		show = shower.Show(key, show)
	}
	cmd(command, obj)

	list := strings.Split(key, ".")
	for _, show := range show {
		key := strings.Split(show.Name, " ")[0]

		log.Debug(key, list[len(list)-1])
		if key == list[len(list)-1] {
			config.Name = show.Name
			config.Help = show.Help
			command.Name = show.Name
			command.Help = show.Help
			continue
		}
		if action, ok := command.Action[key]; ok {
			action.Name = show.Name
			action.Help = show.Help
		}
	}

	last := ice.Index
	for i := 1; i < len(list); i++ {
		has := false
		ice.Pulse.Search(strings.Join(list[:i], ".")+".", func(p *ice.Context, s *ice.Context) {
			has, last = true, s
		})

		if i == len(list)-1 {
			context := &ice.Context{Name: list[len(list)-2],
				Configs: map[string]*ice.Config{list[len(list)-1]: config},
				Commands: map[string]*ice.Command{
					ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
						m.Load()
					}},
					ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
						m.Save()
					}},
					list[len(list)-1]: command,
				},
			}
			log.Debug(fmt.Sprintf("%s %s %s.%s", last.Name, "<-", context.Name))

			if !has {
				last.Register(context, &web.Frame{})
			} else {
				last.Merge(context)
			}
			break
		}

		if !has {
			context := &ice.Context{Name: list[i-1]}
			log.Debug(last.Name, "<-", context.Name)
			last.Register(context, &web.Frame{})
			last = context
		}
	}
}

func _name(key string, index int) string {
	list := strings.Split(key, ".")
	return list[(len(list)+index)%len(list)]
}
