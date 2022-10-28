import * as Blockly from 'blockly/core'
import { Block } from 'blockly/core'
import { javascriptGenerator } from 'blockly/javascript'

// TODO: writing custom function to generate xml is stupid.
// blockly already does it somewhere, figure out how to use that

export const initialXml =
  `<xml xmlns="https://developers.google.com/blockly/xml">
    <block type="advanced_repeat" x="50" y="50">
      <field name="loops">number_of_rotations</field>
      <field name="count">5</field>
      <statement name="inside">
        <block type="character_1">
          <field name="character">Raiden Shogun</field>
          <next>
            <block type="action_attack">
              <field name="count">3</field>
              <next>
                <block type="action_skill">
                  <field name="skill_type"></field>
                  <next>
                    <block type="character_2">
                      <field name="character">Xingqiu</field>
                      <next>
                        <block type="action_attack">
                          <field name="count">1</field>
                        </block>
                      </next>
                    </block>
                  </next>
                </block>
              </next>
            </block>
         </next>
        </block>
      </statement>
    </block>
  </xml>`

export const hotkeyToBlockXml = new Map<string, string>()

// TODO make all characters work
// - grab them from team.tsx or something
// - properly handle iconUrls and 'codeIdentifier'
const characters = ["Raiden Shogun", "Xingqiu", "Bennett", "Xiangling"]

function defineCharacterBlock(index: number) {
  // define block
  const dropdownField = [{
    "type": "field_dropdown",
    "name": "character",
    "options": [
      // fill first option with 'characters[index]'
      [{"src": `https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/${characters[index]}_icon.png`,"width": 45,"height": 45}, characters[index]],
      // fill other option with 'characters' except 'characters[index]' 
      ...characters
        .filter((_, i) => i !== index)
        .map(character => [{"src":`https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/${character}_icon.png`,"width":45,"height":45}, character])
    ]
  }]
  Blockly.Blocks[`character_${index +1}`] = {
    init: function() {
      this.jsonInit({
        "type": `character_${index +1}`,
        "message0": "%1",
        "args0": dropdownField,
        "previousStatement": null,
        "nextStatement": null,
        "colour": "#56458c",
      })
    }
  }
  // define action list code correspoding to the block
  javascriptGenerator["character_" + (index + 1)] = (block: Block) =>
    block.getFieldValue("character").toLowerCase()
      // use first name for multi word character names 
      .split(" ")[0]
      // put space or semicolon
      .concat(block.getNextBlock()?.getFieldValue("character") == null ? " " : ";\n")
  // xml snippet of an example block instance with default values>
  // which should be inserted when editing the xml via hotkeys
  hotkeyToBlockXml.set(
    index + 1 + "", 
    `<next>
      <block type="character_${index + 1}">
        <field name="character">${characters[index]}</field>
      </block>
    </next>`
  )
}

function defineActionBlock(
    blockName: string, 
    assignedHotkeysInEditor: string[], 
    fieldImageSrc: string, 
    BlockMessage: string, 
    codeIdentifier: string, 
    optionalCount?: boolean, 
    optionalDropdown?: { name: string, options: string[][] }, 
    optionalCheckbox?: { name: string, value: string}
) {
  // define block
  const message: string = "%1 " + BlockMessage + [optionalCount, optionalDropdown, optionalCheckbox].filter(x => x !== undefined).map((_, i) => " %" + (i + 2))
  const fields: object[] = [
    {
      "type": "field_image",
      "src": fieldImageSrc, "width": 25, "height": 25,
    }
  ]
  optionalCount !== undefined && fields.push(
     {
      "type": "field_number",
      "name": "count",
      "value": 1, "min": 1, "precision": 1
    }
  )
  optionalDropdown !== undefined && fields.push(
    {
      "type": "field_dropdown",
      "name": optionalDropdown.name,
      "options": optionalDropdown.options
    }
  )
  optionalCheckbox !== undefined && fields.push(
  {
    "type": "field_checkbox",
    "name": optionalCheckbox.name,
    "checked": true
  }
  )
  Blockly.Blocks[blockName] = {
    init: function() {
      this.jsonInit({
        "type": blockName,
        "message0": message,
        "args0": fields,
        "previousStatement": null,
        "nextStatement": null,
        "colour": "#5879a3",
      })
    }
  }
  // define action list code correspoding to the block
  javascriptGenerator[blockName] = (block: Block) => 
    codeIdentifier
      // notate optional action params
      .concat(optionalDropdown ? block.getFieldValue(optionalDropdown.name) : "")
      // notate optional action params
      .concat(optionalCheckbox ? (block.getFieldValue(optionalCheckbox.name)==="TRUE" ? optionalCheckbox.value : "") : "")
      // (maybe) notate optional use count
      .concat(optionalCount ? (block.getFieldValue("count") > 1 ? ":" + block.getFieldValue("count") : "") : "")
      // put comma or semicolon
      // only put comma if there is a next block and it is another action (as identified by the colour of the block)
      .concat(block.getNextBlock()?.getColour() === "#5879a3" ? ", " : ";\n")
  // xml snippet of an example block instance with default values
  // which should be inserted when editing the xml via hotkeys
  const countField = optionalCount ? '<field name="count">1</field>' : ''
  const dropdownField = optionalDropdown ? '<field name="' + optionalDropdown.name + '">' + optionalDropdown.options[0][1] + '</field>' : ''
  const optionField = optionalCheckbox ? '<field name="' + optionalCheckbox.name + '">TRUE</field>' : ''
  assignedHotkeysInEditor.forEach(hotkey => {
    hotkeyToBlockXml.set(
      hotkey,
      `<next>
        <block type="${blockName}">
          ${countField}
          ${dropdownField}
          ${optionField}
        </block>
      </next>`
    )
  })
}

characters.forEach((_, i) => defineCharacterBlock(i))
defineActionBlock(
  "action_burst", 
  ["q", "Q"], 
  "https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/%25E6%2597%2585%25E8%25A1%258C%25E8%2580%2585(%25E9%25A3%258E)/battle_talent_2/battle_talent_2.png",
  "burst",
  "burst",
)
defineActionBlock(
  "action_skill", 
  ["e", "E"],
  "https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/%25E6%2597%2585%25E8%25A1%258C%25E8%2580%2585(%25E9%25A3%258E)/battle_talent_1/battle_talent_1.png",
  "skill",
  "skill",
  undefined,
  { name: "skill_type", options: [["press",""], ["hold","[hold=1]"]] }
)
defineActionBlock(
  "action_attack", 
  ["n", "N"],
  "https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/%25E6%2597%2585%25E8%25A1%258C%25E8%2580%2585(%25E9%25A3%258E)/battle_talent_0/battle_talent_0.png",
  "normal attack", 
  "attack", 
  true
)
// TODO icons from puush links suck, fix that
defineActionBlock(
  "action_charge", 
  ["c", "C"],
  "https://puu.sh/Jq381/58cdc6ef28.png",
  "charged attack", 
  "charge", 
  true
)
if (false) { // TODO can at least 1 party member aim
  defineActionBlock(
    "action_aim", 
    ["a", "A"],
    "https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/%25E7%2594%2598%25E9%259B%25A8/battle_talent_0/battle_talent_0.png",
    "aimed shot", 
    "aim", 
    true,
    { name: "aim_mode", options: [["weakspot", "[weakspot=1]"], ["normal",""]] }
  )
}
if (false) { // TODO can at least 1 party member plunge
  defineActionBlock(
    "action_plunge", 
    ["p", "P"],
    "https://puu.sh/JpWro/1a7cee1c30.png",
    "plunge", 
    "plunge", 
    undefined,
    { name: "plunge_type", options: [["high", "high_plunge"], ["low","low_plunge"]] },
    { name: "collision_mode", value: "[collision=1]" }
  )
}
defineActionBlock(
  "action_dash", 
  ["d", "D", "shift", "Shift"],
  "https://puu.sh/JpWzz/f5cc5881cd.png",
  "dash", 
  "dash"
)
defineActionBlock(
  "action_jump", 
  ["j", "J", " ", "spacebar", "Spacebar"],
  "https://puu.sh/JpWBl/7f03436516.png",
  "jump", 
  "jump"
)
defineActionBlock(
  "action_walk", 
  ["w", "W"], 
  "https://puu.sh/JpWEz/3674e88777.png",
  "walk", 
  "walk"
)

// advanced blocks
Blockly.Blocks["advanced_repeat"] = {
  init: function() {
    this.jsonInit({
      "type": "advanced_repeat",
      "message0": "%1 %2 %3 %4",
      "args0": [
        {
          "type": "field_input",
          "name": "loops",
          "text": "loops"
        },
        {
          "type": "field_number",
          "name": "count",
          "value": 2,
          "min": 1,
          "precision": 1
        },
        {
          "type": "input_dummy"
        },
        {
          "type": "input_statement",
          "name": "inside"
        }
      ],
      "previousStatement": null,
      "nextStatement": null,
      "colour": "#123456"
    })
  }
}
javascriptGenerator["advanced_repeat"] = (block: Block) => { 
  const loops = block.getFieldValue("loops")
  return ""
    // start on new line if its not the first block/line
    .concat(block.getPreviousBlock() ? "\n" : "")
    // while loop
    .concat(
`let ${loops} = ${block.getFieldValue("count")};
while ${loops} {
  ${loops} = ${loops} - 1;
          
          
${javascriptGenerator.statementToCode(block, "inside")}


}
`
)}
Blockly.Blocks["advanced_wait"] = {
  init: function() {
    this.jsonInit({
      "type": "advanced_wait",
      "message0": "wait %1 frames",
      "args0": [
        {
          "type": "field_number",
          "name": "frames",
          "value": 30,
          "min": 1,
          "precision": 1
        }
      ],
      "previousStatement": null,
      "nextStatement": null,
      "colour": "#123456",
    })
  }
}
javascriptGenerator["advanced_wait"] = (block: Block) => "wait(" + block.getFieldValue("frames") + ");\n"

/*
const putSpaceOrSemicolon = (block: Block) => (block.getNextBlock()?.getFieldValue("character") == null) 
  ? " "
  : ";\n"

const putCommaOrSemicolon = (block: Block) => (block.getNextBlock()?.getFieldValue("character") == null) 
  ? ", "
  : ";\n"
  
function maybeNotateUseCount(block: Block) {
  const count = block.getFieldValue("count")
  return (count > 1 ? ":" + count : "")
}
// */

/*
Blockly.Blocks["action_dash"] = {
  init: function() {
    this.appendDummyInput()
        .appendField(new Blockly.FieldImage("https://puu.sh/JpWzz/f5cc5881cd.png", 25, 25))
        .appendField("dash")
        // .appendField(new Blockly.FieldDropdown([["dash", "dash"], ["jump", "jump"], ["walk", "walk"]]), "cancel_method")
    this.setPreviousStatement(true, null)
    this.setNextStatement(true, null)
    this.setColour("#5879a3")
  }
}
javascriptGenerator["action_dash"] = (block: Block) => "dash" + putCommaOrSemicolon(block)
hotkeyToBlockXml.set("d", 
`<next>
  <block type="action_dash">
  </block>
</next>`)
// */

/*
Blockly.Blocks["action_jump"] = {
  init: function() {
    this.appendDummyInput()
      	.appendField(new Blockly.FieldImage("https://puu.sh/JpWBl/7f03436516.png", 25, 25))
        .appendField("jump")
        // .appendField(new Blockly.FieldDropdown([["jump", "jump"], ["dash", "dash"], ["walk", "walk"]]), "cancel_method")
    this.setPreviousStatement(true, null)
    this.setNextStatement(true, null)
    this.setColour("#5879a3")
  }
}
javascriptGenerator["action_jump"] = block => "jump" + putCommaOrSemicolon(block)
exampleActionBlocksXml["action_jump"] = 
`<next>
  <block type="action_jump">
  </block>
</next>`
// */

/*
Blockly.Blocks["action_walk"] = {
  init: function() {
    this.appendDummyInput()
        .appendField(new Blockly.FieldImage("https://puu.sh/JpWEz/3674e88777.png", 25, 25))
        .appendField("walk")
        // .appendField(new Blockly.FieldDropdown([["walk", "walk"], ["dash", "dash"], ["jump", "jump"]]), "cancel_method")
    this.setPreviousStatement(true, null)
    this.setNextStatement(true, null)
    this.setColour("#5879a3")
  }
}
javascriptGenerator["action_walk"] = block => "walk" + putCommaOrSemicolon(block)
exampleActionBlocksXml["action_walk"] = 
`<next>
  <block type="action_walk">
  </block>
</next>`

// javascriptGenerator["action_dash"] = 
// javascriptGenerator["action_jump"] = 
// javascriptGenerator["action_walk"] = block => block.getFieldValue("cancel_method") + putCommaOrSemicolon(block)

// */

/*
Blockly.Blocks["action_burst"] = {
  init: function() {
    this.appendDummyInput()
        .appendField(new Blockly.FieldImage("https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/%25E6%2597%2585%25E8%25A1%258C%25E8%2580%2585(%25E9%25A3%258E)/battle_talent_2/battle_talent_2.png", 25, 25))
        .appendField("burst")
    this.setPreviousStatement(true, null)
    this.setNextStatement(true, null)
    this.setColour("#5879a3")
  }
}
javascriptGenerator["action_burst"] = block => "burst" + putCommaOrSemicolon(block)
exampleActionBlocksXml["action_burst"] = 
`<next>
  <block type="action_burst">
  </block>
</next>`
// */

/*
Blockly.Blocks["action_skill"] = {
  init: function() {
    this.appendDummyInput()
        .appendField(new Blockly.FieldImage("https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/%25E6%2597%2585%25E8%25A1%258C%25E8%2580%2585(%25E9%25A3%258E)/battle_talent_1/battle_talent_1.png", 25, 25))
        .appendField("skill")
    this.appendDummyInput()
        .appendField(new Blockly.FieldDropdown([["press",""], ["hold","[hold=1]"]]), "skill_type")
    this.setInputsInline(true)
    this.setPreviousStatement(true, null)
    this.setNextStatement(true, null)
    this.setColour("#5879a3")
  }
}
javascriptGenerator["action_skill"] = block => "skill" + block.getFieldValue("skill_type") + putCommaOrSemicolon(block)
exampleActionBlocksXml["action_skill"] = 
  `<next>
    <block type="action_skill">
      <field name="skill_type"></field>
    </block>
  </next>`
// */

/*
Blockly.Blocks["action_attack"] = {
  init: function() {
    this.appendDummyInput()
        // .appendField(new Blockly.FieldImage("https://puu.sh/Jq3b7/88eb01de02.png", 25, 25))
        .appendField(new Blockly.FieldImage("https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/%25E6%2597%2585%25E8%25A1%258C%25E8%2580%2585(%25E9%25A3%258E)/battle_talent_0/battle_talent_0.png", 25, 25))
        .appendField("normal attack")
        .appendField(new Blockly.FieldNumber(1, 1, Infinity, 1), "count") // (default, min, max, precision)
      // .appendField(new Blockly.FieldDropdown([["1","1"], ["2","2"], ["3","3"], ["4","4"], ["5","5"], ["6","6"]]), "count")
    this.setPreviousStatement(true, null)
    this.setNextStatement(true, null)
    this.setColour("#5879a3")
  }
}
javascriptGenerator["action_attack"] = (block: Block) => "attack" + maybeNotateUseCount(block) + putCommaOrSemicolon(block)
hotkeyToBlockXml.set("n", 
  `<next>
    <block type="action_attack">
      <field name="count">1</field>
    </block>
  </next>`)
// */

/*
Blockly.Blocks["action_charged"] = {
  init: function() {
    this.appendDummyInput()
        .appendField(new Blockly.FieldImage("https://puu.sh/Jq381/58cdc6ef28.png", 25, 25))
        .appendField("charged attack")
        .appendField(new Blockly.FieldNumber(1, 1, Infinity, 1), "count") // (default, min, max, precision)
    this.setPreviousStatement(true, null)
    this.setNextStatement(true, null)
    this.setColour("#5879a3")
  }
}
javascriptGenerator["action_charged"] = block => "charged" + maybeNotateUseCount(block) + putCommaOrSemicolon(block)
exampleActionBlocksXml["action_charged"] = 
  `<next>
    <block type="action_charged">
      <field name="count">1</field>
    </block>
  </next>`
// */

/*
Blockly.Blocks["action_aim"] = {
  init: function() {
    this.appendDummyInput()
        .appendField(new Blockly.FieldImage("https://wiki.hoyolab.com/_ipx/f_webp/https://bbs.hoyolab.com/hoyowiki/picture/character/%25E7%2594%2598%25E9%259B%25A8/battle_talent_0/battle_talent_0.png", 25, 25))
        // .appendField(new Blockly.FieldImage("https://puu.sh/JpXpu/088aaf4c8a.png", 25, 25))
        .appendField("aimed shot")
        .appendField(new Blockly.FieldNumber(1, 1, Infinity, 1), "count") // (default, min, max, precision)
        .appendField(new Blockly.FieldDropdown([["weakspot", "[weakspot=1]"], ["normal",""]]), "aim_mode")
    this.setPreviousStatement(true, null)
    this.setNextStatement(true, null)
    this.setColour("#5879a3")
  }
}
javascriptGenerator["action_aim"] = block => "aim" + block.getFieldValue("aim_mode") + maybeNotateUseCount(block) + putCommaOrSemicolon(block)
exampleActionBlocksXml["action_aim"] = 
  `<next>
    <block type="action_aim">
      <field name="count">1</field>
      <field name="aim_mode">[weakspot=1]</field>
    </block>
  </next>`
// */

/*
Blockly.Blocks["action_plunge"] = {
  init: function() {
    this.appendDummyInput()
        .appendField(new Blockly.FieldImage("https://puu.sh/JpWro/1a7cee1c30.png", 25, 25)) 
        .appendField("plunge")
        .appendField(new Blockly.FieldDropdown([["high", "high_plunge"], ["low","low_plunge"]]), "plunge_type")
        .appendField("collision")
        .appendField(new Blockly.FieldCheckbox("TRUE"), "collision_mode")
    this.setInputsInline(true)
    this.setPreviousStatement(true, null)
    this.setNextStatement(true, null)
    this.setColour("#5879a3")
  }
}
javascriptGenerator["action_plunge"] = block => block.getFieldValue("plunge_type") + (block.getFieldValue("collision_mode")==="TRUE" ? "[collision=1]" : "") + putCommaOrSemicolon(block)
exampleActionBlocksXml["action_plunge"] = 
  `<next>
    <block type="action_plunge">
      <field name="plunge_type">high_plunge</field>
      <field name="collision_mode">TRUE</field>
    </block>
  </next>`
// */
