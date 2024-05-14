package choose_fonts

import (
	"fmt"

	"kitty/tools/cli"
	"kitty/tools/tty"
	"kitty/tools/tui/loop"
)

var _ = fmt.Print
var debugprintln = tty.DebugPrintln

func main() (rc int, err error) {
	if err = kitty_font_backend.start(); err != nil {
		return 1, err
	}
	defer func() {
		if werr := kitty_font_backend.release(); werr != nil {
			if err == nil {
				err = werr
			}
			if rc == 0 {
				rc = 1
			}
		}
	}()
	lp, err := loop.New()
	if err != nil {
		return 1, err
	}
	lp.MouseTrackingMode(loop.FULL_MOUSE_TRACKING)
	h := &handler{lp: lp}
	lp.OnInitialize = func() (string, error) {
		lp.AllowLineWrapping(false)
		lp.SetWindowTitle(`Choose a font for kitty`)
		return "", h.initialize()
	}
	lp.OnWakeup = h.on_wakeup
	lp.OnEscapeCode = h.on_escape_code
	lp.OnFinalize = func() string {
		h.finalize()
		lp.SetCursorVisible(true)
		return ``
	}
	lp.OnMouseEvent = h.on_mouse_event
	lp.OnResize = func(_, _ loop.ScreenSize) error {
		return h.draw_screen()
	}
	lp.OnKeyEvent = h.on_key_event
	lp.OnText = h.on_text
	err = lp.Run()
	if err != nil {
		return 1, err
	}
	ds := lp.DeathSignalName()
	if ds != "" {
		fmt.Println("Killed by signal: ", ds)
		lp.KillIfSignalled()
		return 1, nil
	}
	return lp.ExitCode(), nil
}

func EntryPoint(root *cli.Command) {
	ans := root.AddSubCommand(&cli.Command{
		Name: "choose-fonts",
		Run: func(cmd *cli.Command, args []string) (rc int, err error) {
			return main()
		},
	})
	clone := root.AddClone(ans.Group, ans)
	clone.Hidden = false
	clone.Name = "choose_fonts"
}