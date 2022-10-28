// TODOs for blockly:
// make all characters work
// - grab them from team.tsx or something
// - properly handle iconUrls and 'codeIdentifier'
// actually give the generated code to the action list config
// persistence 
// - pretty simple, store the xml
// generate xml from action list config file
// - requires full replicability of action list config in UI
// - if not, find workaround for not replicable action list configs
// localization
// - not much to translate, but could look into Blockly.Msg[] and i18n things
// styling
// - fix buggy workspace size, improve layout

// other TODOs:
// add dropdown to select an enemy (lvl and res) + input number of enemies
// could steal from nagi's heroku app or write a script to steal from wiki
// dropdown menu to set active character at the beginning ( for srl <3 )
// ^ these settings should persist locally on browser or for logged in user?
// (same way the other sim settings persist)
// that way, you dont really have to change them/interact with them 
// https://discord.com/channels/845087716541595668/983391844631212112/1035679983131168799


import React from "react"
import Blockly from "blockly"
import { javascriptGenerator } from 'blockly/javascript'
import { hotkeyToBlockXml, initialXml } from "./actionListBlocks"
import { BlocklyWorkspace } from "react-blockly"

const toolboxConfiguration = {
  kind:"flyoutToolbox",
  contents: [
    ...Object.keys(Blockly.Blocks)
      .filter(key => key.includes("character") || key.includes("action") || key.includes("advanced"))
      .map(key => ({kind: "block", type: key}))
  ]
}
const darkTheme = Blockly.Theme.defineTheme("dark", {
    "base": Blockly.Themes.Classic,
    "componentStyles": {
        "workspaceBackgroundColour": "#1E2022",
        "toolboxBackgroundColour": "#313335",
        "toolboxForegroundColour": "#9CA1AD",
        "flyoutBackgroundColour": "#3F4245",
        "flyoutForegroundColour": "#ccc",
        "flyoutOpacity": 0.9,
        "scrollbarColour": "#797979",
        "insertionMarkerColour": "#fff",
        "insertionMarkerOpacity": 0.3,
        "scrollbarOpacity": 0.4,
        "cursorColour": "#d0d0d0",
        "blackBackground": "#333",
    },
})
const workspaceConfiguration = {
  // renderer: "zelos", 
  renderer: "thrasos",
  theme: "dark",
	// horizontalLayout : true, 
	trashcan : true, 
  scrollbars: true,
  zoom: {
    controls : true, 
    wheel : true, 
    startScale : 0.9, 
    maxScale : 1.5, 
    minScale : 0.3, 
    scaleSpeed : 1.05
  }
}

export function ActionListBuilder() {
  const [actionListXml, setActionListXml] = React.useState(initialXml)
  const [actionListCode, setActionListCode] = React.useState("")
  const [areHotkeysEnabled, setAreHotkeysEnabled] = React.useState(false)

  function workspaceDidChange(workspace: any) { // TODO fix 'any'
    setActionListCode(javascriptGenerator.workspaceToCode(workspace))
  }

  function toggleHotkeys() {
    setAreHotkeysEnabled(!areHotkeysEnabled)
  }

  // this will need to be a bit be more sophisticated
  // to handle scopes from loops, ifs, ...
  // could also look into blockly hotkeys 
  function editXmlWithHotkeys(e: React.KeyboardEvent) {
    let newXml = ""
    if (e.key === "Backspace") { // || e.key === "Delete") { 
      const deleteStart = actionListXml.lastIndexOf("<next>"); if (deleteStart === -1) { return }
      const deleteEnd = actionListXml.indexOf("</next>") + "</next>".length; if (deleteEnd === -1) { return }
      newXml = actionListXml.slice(0, deleteStart) + actionListXml.slice(deleteEnd)
    } else {
      if (hotkeyToBlockXml.has(e.key)) {
        const insertedXml = hotkeyToBlockXml.get(e.key)
        const insertPosition = actionListXml.indexOf("</block>"); if (insertPosition === -1) { return }
        newXml = actionListXml.slice(0, insertPosition) + insertedXml + actionListXml.slice(insertPosition)
      } else {
        // pressed key is not a hotkey
        return 
      }
    }
    Blockly.Xml.clearWorkspaceAndLoadFromXml(Blockly.Xml.textToDom(newXml), Blockly.getMainWorkspace())
  }

  return (
    <div>
      <div {...areHotkeysEnabled && {onKeyDown: editXmlWithHotkeys}}>
        <div style={{backgroundColor: "#1e2022", display: "flex"}}>
            <button 
              style={{
                height: "40px", width: "130px", margin: "8px", 
                border: "lightgrey solid 1px", borderRadius: "3px", 
                color: "black", backgroundColor: areHotkeysEnabled ? "#abc" : "#ddd"
              }} 
              onClick={toggleHotkeys}
            >
              {areHotkeysEnabled ? "disable" : "enable"} hotkeys
            </button>
            {areHotkeysEnabled && <p style={{color: "white"}}>hotkeys {[...hotkeyToBlockXml.keys()].toString()} are enabled. if they dont work, click anywhere inside the workspace window to re-enable them or turn them off and on again</p>}
        </div>
        <BlocklyWorkspace
          toolboxConfiguration={toolboxConfiguration}
          workspaceConfiguration={workspaceConfiguration}
          initialXml={actionListXml}
          onXmlChange={setActionListXml}
          onWorkspaceChange={workspaceDidChange}
          className="blockly-workspace"
        />
        <pre
        //  id="generated-xml"
        >{actionListXml}</pre>
        <textarea
          // id="code"
          style={{height: "800px", width: "800px", fontFamily: "monospace", padding: "3px"}}
          value={actionListCode}
          readOnly
        ></textarea>
      </div>
    </div>
  )
}
