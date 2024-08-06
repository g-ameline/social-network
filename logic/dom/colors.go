package dom

import (
	"fmt"
)

func To_color(name string, shade func(int) string) string {
	var N, n int
	for _, byte := range []byte(name) {
		integer := int(byte)
		n = -n + integer
		N += ((integer * 811) * (integer * 811)) % 360
	}
	if n < 0 {
		n = -n
	}
	return shade(N + n)
}

func (the_node Node) Outline(name string) Node {
	return the_node.
		Style(Color, To_color(name, Deep)).
		Style(Background_color, Ivory).
		Style(Border_style, "solid").
		Style(Border_width, "thick").
		Style(Border_radius, "10px").
		Style(Font_weight, "800").
		Style(Box_sizing, Border_box).
		Style(Border_color, To_color(name, Deep))
}
func (the_node Node) Outline_Rec(name string) Node {
	the_node.Outline(name)
	for _, kid_node := range the_node.Get_kids() {
		kid_node.Outline_Rec(name)
	}
	return the_node
}

func (the_node Node) Tinge(name string) Node {
	return the_node.
		Style(Background_color, To_color(name, Light)).
		Style(Color, To_color(name, Dark))
}

func (the_node Node) Murk(name string) Node {
	return the_node.
		Style(Background_color, To_color(name, Dark)).
		Style(Color, To_color(name, Light))
}
func (the_node Node) Circle(name string) Node {
	return the_node.
		Style(Background_color, To_color(name, Deep)).
		Style(Color, "black").
		Style(Border_width, "medium").
		Style(Border_color, "black").
		Style(Border_style, "solid").
		Style(Border_radius, "5px").
		Style(Font_style, "italic").
		Style(Font_size, "20px").
		Style(Text_align, "center").
		Style(Padding, "8px").
		Style(Box_sizing, Border_box)
}

func (the_node Node) Circle_Rec(name string) Node {
	the_node.Circle(name)
	for _, kid_node := range the_node.Get_kids() {
		kid_node.Circle_Rec(name)
	}
	return the_node
}

func HSL(H, S, L int) string {
	return fmt.Sprintf("hsl(%v, %v%%, %v%%)", H%360, S%100, L%100)
}

func Light(H int) string {
	return fmt.Sprintf("hsl(%v, %v%%, %v%%)", H%360, 100, 80)
}

func Dark(H int) string {
	return fmt.Sprintf("hsl(%v, %v%%, %v%%)", H%360, 100, 20)
}

func Vivid(H int) string {
	return fmt.Sprintf("hsl(%v, %v%%, %v%%)", H%360, 100, 50)
}

func Deep(H int) string {
	return fmt.Sprintf("hsl(%v, %v%%, %v%%)", H%360, 100, 40)
}

func (the_node Node) Text_black() Node {
	return the_node.Style(Color, "black")
}
func (the_node Node) Indian_red() Node {
	return the_node.Style(Background_color, Indian_red)
}
func (the_node Node) Light_coral() Node {
	return the_node.Style(Background_color, Light_coral)
}
func (the_node Node) Salmon() Node {
	return the_node.Style(Background_color, Salmon)
}
func (the_node Node) Dark_salmon() Node {
	return the_node.Style(Background_color, Dark_salmon)
}
func (the_node Node) Light_salmon() Node {
	return the_node.Style(Background_color, Light_salmon)
}
func (the_node Node) Crimson() Node {
	return the_node.Style(Background_color, Crimson)
}
func (the_node Node) Red() Node {
	return the_node.Style(Background_color, Red)
}
func (the_node Node) Fire_brick() Node {
	return the_node.Style(Background_color, Fire_brick)
}
func (the_node Node) Dark_red() Node {
	return the_node.Style(Background_color, Dark_red)
}
func (the_node Node) Pink() Node {
	return the_node.Style(Background_color, Pink)
}
func (the_node Node) Light_pink() Node {
	return the_node.Style(Background_color, Light_pink)
}
func (the_node Node) Hot_pink() Node {
	return the_node.Style(Background_color, Hot_pink)
}
func (the_node Node) Deep_pink() Node {
	return the_node.Style(Background_color, Deep_pink)
}
func (the_node Node) Medium_violet_red() Node {
	return the_node.Style(Background_color, Medium_violet_red)
}
func (the_node Node) Pale_violet_red() Node {
	return the_node.Style(Background_color, Pale_violet_red)
}
func (the_node Node) Orange() Node {
	return the_node.Style(Background_color, Orange)
}
func (the_node Node) Coral() Node {
	return the_node.Style(Background_color, Coral)
}
func (the_node Node) Tomato() Node {
	return the_node.Style(Background_color, Tomato)
}
func (the_node Node) Orange_red() Node {
	return the_node.Style(Background_color, Orange_red)
}
func (the_node Node) Dark_orange() Node {
	return the_node.Style(Background_color, Dark_orange)
}
func (the_node Node) Yellow() Node {
	return the_node.Style(Background_color, Yellow)
}
func (the_node Node) Gold() Node {
	return the_node.Style(Background_color, Gold)
}
func (the_node Node) Light_yellow() Node {
	return the_node.Style(Background_color, Light_yellow)
}
func (the_node Node) Lemon_chiffon() Node {
	return the_node.Style(Background_color, Lemon_chiffon)
}
func (the_node Node) Light_goldenrod_yellow() Node {
	return the_node.Style(Background_color, Light_goldenrod_yellow)
}
func (the_node Node) Papaya_whip() Node {
	return the_node.Style(Background_color, Papaya_whip)
}
func (the_node Node) Moccasin() Node {
	return the_node.Style(Background_color, Moccasin)
}
func (the_node Node) Peach_puff() Node {
	return the_node.Style(Background_color, Peach_puff)
}
func (the_node Node) Pale_goldenrod() Node {
	return the_node.Style(Background_color, Pale_goldenrod)
}
func (the_node Node) Khaki() Node {
	return the_node.Style(Background_color, Khaki)
}
func (the_node Node) Dark_khaki() Node {
	return the_node.Style(Background_color, Dark_khaki)
}
func (the_node Node) Purple() Node {
	return the_node.Style(Background_color, Purple)
}
func (the_node Node) Lavender() Node {
	return the_node.Style(Background_color, Lavender)
}
func (the_node Node) Thistle() Node {
	return the_node.Style(Background_color, Thistle)
}
func (the_node Node) Plum() Node {
	return the_node.Style(Background_color, Plum)
}
func (the_node Node) Violet() Node {
	return the_node.Style(Background_color, Violet)
}
func (the_node Node) Orchid() Node {
	return the_node.Style(Background_color, Orchid)
}
func (the_node Node) Fuchsia() Node {
	return the_node.Style(Background_color, Fuchsia)
}
func (the_node Node) Magenta() Node {
	return the_node.Style(Background_color, Magenta)
}
func (the_node Node) Medium_orchid() Node {
	return the_node.Style(Background_color, Medium_orchid)
}
func (the_node Node) Medium_purple() Node {
	return the_node.Style(Background_color, Medium_purple)
}
func (the_node Node) Rebecca_purple() Node {
	return the_node.Style(Background_color, Rebecca_purple)
}
func (the_node Node) Blue_violet() Node {
	return the_node.Style(Background_color, Blue_violet)
}
func (the_node Node) Dark_violet() Node {
	return the_node.Style(Background_color, Dark_violet)
}
func (the_node Node) Dark_orchid() Node {
	return the_node.Style(Background_color, Dark_orchid)
}
func (the_node Node) Dark_magenta() Node {
	return the_node.Style(Background_color, Dark_magenta)
}
func (the_node Node) Indigo() Node {
	return the_node.Style(Background_color, Indigo)
}
func (the_node Node) Slate_blue() Node {
	return the_node.Style(Background_color, Slate_blue)
}
func (the_node Node) Dark_slate_blue() Node {
	return the_node.Style(Background_color, Dark_slate_blue)
}
func (the_node Node) Medium_slate_blue() Node {
	return the_node.Style(Background_color, Medium_slate_blue)
}
func (the_node Node) Green() Node {
	return the_node.Style(Background_color, Green)
}
func (the_node Node) Green_yellow() Node {
	return the_node.Style(Background_color, Green_yellow)
}
func (the_node Node) Chartreuse() Node {
	return the_node.Style(Background_color, Chartreuse)
}
func (the_node Node) Lawn_green() Node {
	return the_node.Style(Background_color, Lawn_green)
}
func (the_node Node) Lime() Node {
	return the_node.Style(Background_color, Lime)
}
func (the_node Node) Lime_green() Node {
	return the_node.Style(Background_color, Lime_green)
}
func (the_node Node) Pale_green() Node {
	return the_node.Style(Background_color, Pale_green)
}
func (the_node Node) Light_green() Node {
	return the_node.Style(Background_color, Light_green)
}
func (the_node Node) Medium_spring_green() Node {
	return the_node.Style(Background_color, Medium_spring_green)
}
func (the_node Node) Spring_green() Node {
	return the_node.Style(Background_color, Spring_green)
}
func (the_node Node) Medium_sea_green() Node {
	return the_node.Style(Background_color, Medium_sea_green)
}
func (the_node Node) Sea_green() Node {
	return the_node.Style(Background_color, Sea_green)
}
func (the_node Node) Forest_green() Node {
	return the_node.Style(Background_color, Forest_green)
}
func (the_node Node) Dark_green() Node {
	return the_node.Style(Background_color, Dark_green)
}
func (the_node Node) Yellow_green() Node {
	return the_node.Style(Background_color, Yellow_green)
}
func (the_node Node) Olive_drab() Node {
	return the_node.Style(Background_color, Olive_drab)
}
func (the_node Node) Olive() Node {
	return the_node.Style(Background_color, Olive)
}
func (the_node Node) Dark_olive_green() Node {
	return the_node.Style(Background_color, Dark_olive_green)
}
func (the_node Node) Medium_aquamarine() Node {
	return the_node.Style(Background_color, Medium_aquamarine)
}
func (the_node Node) Dark_sea_green() Node {
	return the_node.Style(Background_color, Dark_sea_green)
}
func (the_node Node) Light_sea_green() Node {
	return the_node.Style(Background_color, Light_sea_green)
}
func (the_node Node) Dark_cyan() Node {
	return the_node.Style(Background_color, Dark_cyan)
}
func (the_node Node) Teal() Node {
	return the_node.Style(Background_color, Teal)
}
func (the_node Node) Blue() Node {
	return the_node.Style(Background_color, Blue)
}
func (the_node Node) Aqua() Node {
	return the_node.Style(Background_color, Aqua)
}
func (the_node Node) Cyan() Node {
	return the_node.Style(Background_color, Cyan)
}
func (the_node Node) Light_cyan() Node {
	return the_node.Style(Background_color, Light_cyan)
}
func (the_node Node) Pale_turquoise() Node {
	return the_node.Style(Background_color, Pale_turquoise)
}
func (the_node Node) Aquamarine() Node {
	return the_node.Style(Background_color, Aquamarine)
}
func (the_node Node) Turquoise() Node {
	return the_node.Style(Background_color, Turquoise)
}
func (the_node Node) Medium_turquoise() Node {
	return the_node.Style(Background_color, Medium_turquoise)
}
func (the_node Node) Dark_turquoise() Node {
	return the_node.Style(Background_color, Dark_turquoise)
}
func (the_node Node) Cadet_blue() Node {
	return the_node.Style(Background_color, Cadet_blue)
}
func (the_node Node) Steel_blue() Node {
	return the_node.Style(Background_color, Steel_blue)
}
func (the_node Node) Light_steel_blue() Node {
	return the_node.Style(Background_color, Light_steel_blue)
}
func (the_node Node) Powder_blue() Node {
	return the_node.Style(Background_color, Powder_blue)
}
func (the_node Node) Light_blue() Node {
	return the_node.Style(Background_color, Light_blue)
}
func (the_node Node) Sky_blue() Node {
	return the_node.Style(Background_color, Sky_blue)
}
func (the_node Node) Light_sky_blue() Node {
	return the_node.Style(Background_color, Light_sky_blue)
}
func (the_node Node) Deep_sky_blue() Node {
	return the_node.Style(Background_color, Deep_sky_blue)
}
func (the_node Node) Dodger_blue() Node {
	return the_node.Style(Background_color, Dodger_blue)
}
func (the_node Node) Cornflower_blue() Node {
	return the_node.Style(Background_color, Cornflower_blue)
}
func (the_node Node) Royal_blue() Node {
	return the_node.Style(Background_color, Royal_blue)
}
func (the_node Node) Medium_blue() Node {
	return the_node.Style(Background_color, Medium_blue)
}
func (the_node Node) Dark_blue() Node {
	return the_node.Style(Background_color, Dark_blue)
}
func (the_node Node) Navy() Node {
	return the_node.Style(Background_color, Navy)
}
func (the_node Node) Midnight_blue() Node {
	return the_node.Style(Background_color, Midnight_blue)
}
func (the_node Node) Brown() Node {
	return the_node.Style(Background_color, Brown)
}
func (the_node Node) Cornsilk() Node {
	return the_node.Style(Background_color, Cornsilk)
}
func (the_node Node) Blanched_almond() Node {
	return the_node.Style(Background_color, Blanched_almond)
}
func (the_node Node) Bisque() Node {
	return the_node.Style(Background_color, Bisque)
}
func (the_node Node) Navajo_white() Node {
	return the_node.Style(Background_color, Navajo_white)
}
func (the_node Node) Wheat() Node {
	return the_node.Style(Background_color, Wheat)
}
func (the_node Node) Burly_wood() Node {
	return the_node.Style(Background_color, Burly_wood)
}
func (the_node Node) Tan() Node {
	return the_node.Style(Background_color, Tan)
}
func (the_node Node) Rosy_brown() Node {
	return the_node.Style(Background_color, Rosy_brown)
}
func (the_node Node) Sandy_brown() Node {
	return the_node.Style(Background_color, Sandy_brown)
}
func (the_node Node) Goldenrod() Node {
	return the_node.Style(Background_color, Goldenrod)
}
func (the_node Node) Dark_goldenrod() Node {
	return the_node.Style(Background_color, Dark_goldenrod)
}
func (the_node Node) Peru() Node {
	return the_node.Style(Background_color, Peru)
}
func (the_node Node) Chocolate() Node {
	return the_node.Style(Background_color, Chocolate)
}
func (the_node Node) Saddle_brown() Node {
	return the_node.Style(Background_color, Saddle_brown)
}
func (the_node Node) Sienna() Node {
	return the_node.Style(Background_color, Sienna)
}
func (the_node Node) Maroon() Node {
	return the_node.Style(Background_color, Maroon)
}
func (the_node Node) White() Node {
	return the_node.Style(Background_color, White)
}
func (the_node Node) Snow() Node {
	return the_node.Style(Background_color, Snow)
}
func (the_node Node) Honey_dew() Node {
	return the_node.Style(Background_color, Honey_dew)
}
func (the_node Node) Mint_cream() Node {
	return the_node.Style(Background_color, Mint_cream)
}
func (the_node Node) Azure() Node {
	return the_node.Style(Background_color, Azure)
}
func (the_node Node) Alice_blue() Node {
	return the_node.Style(Background_color, Alice_blue)
}
func (the_node Node) Ghost_white() Node {
	return the_node.Style(Background_color, Ghost_white)
}
func (the_node Node) White_smoke() Node {
	return the_node.Style(Background_color, White_smoke)
}
func (the_node Node) Sea_shell() Node {
	return the_node.Style(Background_color, Sea_shell)
}
func (the_node Node) Beige() Node {
	return the_node.Style(Background_color, Beige)
}
func (the_node Node) Old_lace() Node {
	return the_node.Style(Background_color, Old_lace)
}
func (the_node Node) Floral_white() Node {
	return the_node.Style(Background_color, Floral_white)
}
func (the_node Node) Ivory() Node {
	return the_node.Style(Background_color, Ivory)
}
func (the_node Node) Antique_white() Node {
	return the_node.Style(Background_color, Antique_white)
}
func (the_node Node) Linen() Node {
	return the_node.Style(Background_color, Linen)
}
func (the_node Node) Lavender_blush() Node {
	return the_node.Style(Background_color, Lavender_blush)
}
func (the_node Node) Misty_rose() Node {
	return the_node.Style(Background_color, Misty_rose)
}
func (the_node Node) Gray() Node {
	return the_node.Style(Background_color, Gray)
}
func (the_node Node) Gainsboro() Node {
	return the_node.Style(Background_color, Gainsboro)
}
func (the_node Node) Light_gray() Node {
	return the_node.Style(Background_color, Light_gray)
}
func (the_node Node) Silver() Node {
	return the_node.Style(Background_color, Silver)
}
func (the_node Node) Dark_gray() Node {
	return the_node.Style(Background_color, Dark_gray)
}
func (the_node Node) Dim_gray() Node {
	return the_node.Style(Background_color, Dim_gray)
}
func (the_node Node) Light_slate_gray() Node {
	return the_node.Style(Background_color, Light_slate_gray)
}
func (the_node Node) Slate_gray() Node {
	return the_node.Style(Background_color, Slate_gray)
}
func (the_node Node) Dark_slate_gray() Node {
	return the_node.Style(Background_color, Dark_slate_gray)
}

const (
	Indian_red             = "IndianRed"
	Light_coral            = "LightCoral"
	Salmon                 = "Salmon"
	Dark_salmon            = "DarkSalmon"
	Light_salmon           = "LightSalmon"
	Crimson                = "Crimson"
	Red                    = "Red"
	Fire_brick             = "FireBrick"
	Dark_red               = "DarkRed"
	Pink                   = "Pink"
	Light_pink             = "LightPink"
	Hot_pink               = "HotPink"
	Deep_pink              = "DeepPink"
	Medium_violet_red      = "MediumVioletRed"
	Pale_violet_red        = "PaleVioletRed"
	Orange                 = "Orange"
	COLOR                  = "COLOR"
	Coral                  = "Coral"
	Tomato                 = "Tomato"
	Orange_red             = "OrangeRed"
	Dark_orange            = "DarkOrange"
	Yellow                 = "Yellow"
	Gold                   = "Gold"
	Light_yellow           = "LightYellow"
	Lemon_chiffon          = "LemonChiffon"
	Light_goldenrod_yellow = "LightGoldenrodYellow"
	Papaya_whip            = "PapayaWhip"
	Moccasin               = "Moccasin"
	Peach_puff             = "PeachPuff"
	Pale_goldenrod         = "PaleGoldenrod"
	Khaki                  = "Khaki"
	Dark_khaki             = "DarkKhaki"
	Purple                 = "Purple"
	Lavender               = "Lavender"
	Thistle                = "Thistle"
	Plum                   = "Plum"
	Violet                 = "Violet"
	Orchid                 = "Orchid"
	Fuchsia                = "Fuchsia"
	Magenta                = "Magenta"
	Medium_orchid          = "MediumOrchid"
	Medium_purple          = "MediumPurple"
	Rebecca_purple         = "RebeccaPurple"
	Blue_violet            = "BlueViolet"
	Dark_violet            = "DarkViolet"
	Dark_orchid            = "DarkOrchid"
	Dark_magenta           = "DarkMagenta"
	Indigo                 = "Indigo"
	Slate_blue             = "SlateBlue"
	Dark_slate_blue        = "DarkSlateBlue"
	Medium_slate_blue      = "MediumSlateBlue"
	Green                  = "Green"
	Green_yellow           = "GreenYellow"
	Chartreuse             = "Chartreuse"
	Lawn_green             = "LawnGreen"
	Lime                   = "Lime"
	Lime_green             = "LimeGreen"
	Pale_green             = "PaleGreen"
	Light_green            = "LightGreen"
	Medium_spring_green    = "MediumSpringGreen"
	Spring_green           = "SpringGreen"
	Medium_sea_green       = "MediumSeaGreen"
	Sea_green              = "SeaGreen"
	Forest_green           = "ForestGreen"
	Dark_green             = "DarkGreen"
	Yellow_green           = "YellowGreen"
	Olive_drab             = "OliveDrab"
	Olive                  = "Olive"
	Dark_olive_green       = "DarkOliveGreen"
	Medium_aquamarine      = "MediumAquamarine"
	Dark_sea_green         = "DarkSeaGreen"
	Light_sea_green        = "LightSeaGreen"
	Dark_cyan              = "DarkCyan"
	Teal                   = "Teal"
	Blue                   = "Blue"
	Aqua                   = "Aqua"
	Cyan                   = "Cyan"
	Light_cyan             = "LightCyan"
	Pale_turquoise         = "PaleTurquoise"
	Aquamarine             = "Aquamarine"
	Turquoise              = "Turquoise"
	Medium_turquoise       = "MediumTurquoise"
	Dark_turquoise         = "DarkTurquoise"
	Cadet_blue             = "CadetBlue"
	Steel_blue             = "SteelBlue"
	Light_steel_blue       = "LightSteelBlue"
	Powder_blue            = "PowderBlue"
	Light_blue             = "LightBlue"
	Sky_blue               = "SkyBlue"
	Light_sky_blue         = "LightSkyBlue"
	Deep_sky_blue          = "DeepSkyBlue"
	Dodger_blue            = "DodgerBlue"
	Cornflower_blue        = "CornflowerBlue"
	Royal_blue             = "RoyalBlue"
	Medium_blue            = "MediumBlue"
	Dark_blue              = "DarkBlue"
	Navy                   = "Navy"
	Midnight_blue          = "MidnightBlue"
	Brown                  = "Brown"
	Cornsilk               = "Cornsilk"
	Blanched_almond        = "BlanchedAlmond"
	Bisque                 = "Bisque"
	Navajo_white           = "NavajoWhite"
	Wheat                  = "Wheat"
	Burly_wood             = "BurlyWood"
	Tan                    = "Tan"
	Rosy_brown             = "RosyBrown"
	Sandy_brown            = "SandyBrown"
	Goldenrod              = "Goldenrod"
	Dark_goldenrod         = "DarkGoldenrod"
	Peru                   = "Peru"
	Chocolate              = "Chocolate"
	Saddle_brown           = "SaddleBrown"
	Sienna                 = "Sienna"
	Maroon                 = "Maroon"
	White                  = "White"
	Snow                   = "Snow"
	Honey_dew              = "HoneyDew"
	Mint_cream             = "MintCream"
	Azure                  = "Azure"
	Alice_blue             = "AliceBlue"
	Ghost_white            = "GhostWhite"
	White_smoke            = "WhiteSmoke"
	Sea_shell              = "SeaShell"
	Beige                  = "Beige"
	Old_lace               = "OldLace"
	Floral_white           = "FloralWhite"
	Ivory                  = "Ivory"
	Antique_white          = "AntiqueWhite"
	Linen                  = "Linen"
	Lavender_blush         = "LavenderBlush"
	Misty_rose             = "MistyRose"
	Gray                   = "Gray"
	Gainsboro              = "Gainsboro"
	Light_gray             = "LightGray"
	Silver                 = "Silver"
	Dark_gray              = "DarkGray"
	Dim_gray               = "DimGray"
	Light_slate_gray       = "LightSlateGray"
	Slate_gray             = "SlateGray"
	Dark_slate_gray        = "DarkSlateGray"
)
