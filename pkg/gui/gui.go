package gui

import (
	"github.com/Recruit-CSIRT/macApfsMounter/pkg/conf"
	"github.com/Recruit-CSIRT/macApfsMounter/pkg/mt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func Run(config *conf.Config) {
	app := widgets.NewQApplication(0, nil)

	window := createWindow(config)

	window.Show()

	app.Exec()
}

func createWindow(config *conf.Config) *widgets.QMainWindow {

	// make window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(400, 200)
	window.SetWindowTitle("macOS APFS Mounter")

	// gui group of the e01 path
	var (
		imgPathGroup    = widgets.NewQGroupBox2("Evidence path setting", nil)
		imgPathLabel = widgets.NewQLabel2("Image Path", nil, 0)
		imgPathLineEdit = widgets.NewQLineEdit2("/", nil)
		imgPathButton   = widgets.NewQPushButton2("...", nil)
	)
	imgPathButton.ConnectClicked(func(bool) {
		openFileDialog(imgPathLineEdit.Text(), imgPathLineEdit)
	})

	// gui group of filevault password
	var (
		optionsGroup    = widgets.NewQGroupBox2("Options", nil)

		filevaultPassLabel    = widgets.NewQLabel2("FileVault Password", nil, 0)
		filevaultPassLineEdit = widgets.NewQLineEdit2("", nil)

		filetypeLabel    = widgets.NewQLabel2("FileType: ", nil, 0)

		filetypeGroupBox = widgets.NewQGroupBox(nil)
		filetypeBoxLayout = widgets.NewQHBoxLayout()

		filetypeButtonGroup = widgets.NewQButtonGroup(nil)
		filetypeRadioEwf = widgets.NewQRadioButton2("E01", nil)
		//filetypeRadioRaw = widgets.NewQRadioButton2("raw", nil)
		filetypeRadioDmg = widgets.NewQRadioButton2("dmg", nil)
	)

	// conf.SupportFileTypeと一致している必要あり
	filetypeButtonGroup.AddButton(filetypeRadioEwf, 0)
	//filetypeButtonGroup.AddButton(filetypeRadioRaw, 1)
	filetypeButtonGroup.AddButton(filetypeRadioDmg, 2)

	filetypeRadioEwf.SetChecked(true)

	filetypeBoxLayout.AddWidget(filetypeRadioEwf, 0, 0)
	//filetypeBoxLayout.AddWidget(filetypeRadioRaw, 0, 0)
	filetypeBoxLayout.AddWidget(filetypeRadioDmg, 0, 0)
	filetypeGroupBox.SetLayout(filetypeBoxLayout)

	// run button and action
	runButton := widgets.NewQPushButton2("Run", nil)
	runButton.ConnectClicked(func(bool) {

		// get fields value
		config.ImgPath = imgPathLineEdit.Text()
		config.VaultPW = filevaultPassLineEdit.Text()

		config.FileType = conf.SupportFileType[filetypeButtonGroup.CheckedId()]

		if f, err := config.CheckImgFile(); !f {
			widgets.QMessageBox_Information(nil, "Status", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			return
		}

		if err := mt.Run(config); err != nil {
			message := "Failed to mount: " + err.Error()
			widgets.QMessageBox_Information(nil, "Status", message, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else {
			widgets.QMessageBox_Information(nil, "Status", "Finished!", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}

	})

	/*
	closeButton := widgets.NewQPushButton2("Close", nil)
	closeButton.ConnectClicked(func(bool) {
		window.Close()
	})
	*/
	// layout of img path
	var imgPathLayout = widgets.NewQGridLayout2()
	imgPathLayout.AddWidget2(imgPathLabel, 0, 0, 0)
	imgPathLayout.AddWidget2(imgPathLineEdit, 1, 0, 0)
	imgPathLayout.AddWidget2(imgPathButton, 1, 1, 0)
	imgPathGroup.SetLayout(imgPathLayout)

	// layout of option
	var optionsLayout = widgets.NewQGridLayout2()
	optionsLayout.AddWidget2(filevaultPassLabel, 0, 0, core.Qt__AlignLeft)
	optionsLayout.AddWidget2(filevaultPassLineEdit, 0, 1,0)

	optionsLayout.AddWidget2(filetypeLabel, 1, 0 , core.Qt__AlignLeft)
	optionsLayout.AddWidget2(filetypeGroupBox, 1, 1, core.Qt__AlignLeft)
	optionsGroup.SetLayout(optionsLayout)

	// layout setting
	var layout = widgets.NewQGridLayout2()
	layout.AddWidget3(imgPathGroup, 0,0,1,2, 0)
	layout.AddWidget3(optionsGroup, 1,0,1,2, 0)

	//layout.AddWidget3(closeButton, 9, 0, 1, 1, core.Qt__AlignRight)
	layout.AddWidget3(runButton, 9, 1, 1, 1, core.Qt__AlignRight)

	// set layout
	var widget = widgets.NewQWidget(window, 0)
	widget.SetLayout(layout)
	window.SetCentralWidget(widget)

	return window
}


func openFileDialog(path string, lineEdit *widgets.QLineEdit) {

	fileDialog := widgets.NewQFileDialog2(nil, "Path", path, "")
	if fileDialog.Exec() != int(widgets.QDialog__Accepted) {
		return
	}

	selectedFilePath := fileDialog.SelectedFiles()[0]

	lineEdit.SetText(selectedFilePath)
}
