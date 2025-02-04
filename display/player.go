package display

import (
    "strconv"

    "github.com/Ariemeth/frame_assault/mech"
    tl "github.com/Ariemeth/termloop"
)

const (
    textLineStartX = 1    // X offset for text from display edge
    textLineStartY = 1    // Y offset for first text line
    textLineSpacing = 1   // Spacing between text lines
    displayWidth = 25     // Width of the status display
    displayHeight = 12    // Height of the status display (9 text lines + margins)
    numTextLines = 9      // Total number of text lines in display
)

//Player represents a player status display
type Player struct {
    Status
    player      *mech.PlayerMech
    timeSystem  TimeSystemInterface
    textLine1   *tl.Text
    textLine2   *tl.Text
    textLine3   *tl.Text
    textLine4   *tl.Text
    textLine5   *tl.Text
    textLine6   *tl.Text
    textLine7   *tl.Text
    textLine8   *tl.Text
    textLine9   *tl.Text
}

// TimeSystemInterface defines the methods required for time display
type TimeSystemInterface interface {
    FormatGameTime() string
}

//NewPlayer creates a new status display for the specified PlayerMech
func NewPlayer(x, y int, player *mech.PlayerMech, timeSystem TimeSystemInterface, level *tl.BaseLevel) *Player {
    display := &Player{
        Status:     *NewStatus(x, y, displayWidth, displayHeight, level),
        player:     player,
        timeSystem: timeSystem,
        textLine1:  tl.NewText(x, y, "", tl.ColorWhite, tl.ColorBlack),
        textLine2:  tl.NewText(x, y+1, "", tl.ColorWhite, tl.ColorBlack),
        textLine3:  tl.NewText(x, y+2, "", tl.ColorWhite, tl.ColorBlack),
        textLine4:  tl.NewText(x, y+3, "", tl.ColorWhite, tl.ColorBlack),
        textLine5:  tl.NewText(x, y+4, "", tl.ColorWhite, tl.ColorBlack),
        textLine6:  tl.NewText(x, y+5, "", tl.ColorWhite, tl.ColorBlack),
        textLine7:  tl.NewText(x, y+6, "", tl.ColorWhite, tl.ColorBlack),
        textLine8:  tl.NewText(x, y+7, "", tl.ColorWhite, tl.ColorBlack),
        textLine9:  tl.NewText(x, y+8, "", tl.ColorWhite, tl.ColorBlack),
    }
    return display
}

// positionTextLines updates the position of all text lines based on the current offset
func (display *Player) positionTextLines(offsetX, offsetY int) {
    lines := []*tl.Text{
        display.textLine1, display.textLine2, display.textLine3,
        display.textLine4, display.textLine5, display.textLine6,
        display.textLine7, display.textLine8, display.textLine9,
    }
    
    for i, line := range lines {
        x := -offsetX + display.x + textLineStartX
        y := -offsetY + display.y + textLineStartY + (i * textLineSpacing)
        line.SetPosition(x, y)
    }
}

// drawTextLines draws all text lines to the screen
func (display *Player) drawTextLines(screen *tl.Screen) {
    lines := []*tl.Text{
        display.textLine1, display.textLine2, display.textLine3,
        display.textLine4, display.textLine5, display.textLine6,
        display.textLine7, display.textLine8, display.textLine9,
    }
    
    for _, line := range lines {
        line.Draw(screen)
    }
}

// Draw passes the draw call to entity.
func (display *Player) Draw(screen *tl.Screen) {
    offSetX, offSetY := display.level.Offset()
    
    // Draw background
    display.background.SetPosition(-offSetX+display.x, -offSetY+display.y)
    display.background.Draw(screen)
    
    // Position and draw text lines
    display.positionTextLines(offSetX, offSetY)
    display.drawTextLines(screen)
}

// Tick is called to process 1 tick of actions based on the
// current state of the game.
func (display *Player) Tick(event tl.Event) {
    // Time display at the top
    if display.timeSystem != nil {
        display.textLine1.SetText(display.timeSystem.FormatGameTime())
    }
    
    // Player info moved down one line
    display.textLine2.SetText(display.player.Name())
    display.textLine3.SetText("Struture: " + strconv.Itoa(display.player.StructureLeft()))
    x, y := display.player.Position()
    display.textLine4.SetText("Location: (" + strconv.Itoa(x) + "," + strconv.Itoa(y) + ")")

    //assume for now there is only 1 Weapon
    display.textLine5.SetText("Weapons")
    weapons := display.player.Weapons()
    if len(weapons) > 0 {
        display.textLine6.SetText("    Name: " + weapons[0].Name())
        display.textLine6.SetColor(tl.ColorWhite, tl.ColorBlack)
        display.textLine7.SetText("   Range: " + strconv.Itoa(weapons[0].Range()))
        display.textLine8.SetText("  Damage: " + strconv.Itoa(weapons[0].Damage()))
        display.textLine9.SetText("Accuracy: " + strconv.FormatFloat(weapons[0].Accuracy()*100, 'f', 1, 64) + "%")
    } else {
        display.textLine6.SetText("    None")
        display.textLine6.SetColor(tl.ColorRed, tl.ColorBlack)
        display.textLine7.SetText("")
        display.textLine8.SetText("")
        display.textLine9.SetText("")
    }
}
