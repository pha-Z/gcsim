import React from "react"
import { Button, Icon } from '@blueprintjs/core'

// TODO get proper possible actions (which are also based on character?)
// likely needs to be changed entirely.
// Look into actions and characters from 'core', 'simActions' or whatever the thing is called
const availableActions: string[] = [
    "normal attack",
    "charge attack",
    "skill",
    "burst",
    "dash",
    "jump",
]
const availableCharacters: string[] = ["raiden", "ayaka"]

// TODO at hint descriptions for elements
export function ActionListBuilder() {
    const [characterActions, setCharacterActions] = React.useState(
        [
            {character: "raiden", actions: ["normal attack", "skill"]}, 
            {character: "ayaka", actions: ["dash"]}
        ]
    )
    
    // for debug
    React.useEffect(() => console.log(characterActions), [characterActions])

    function handleAddCharacterAction(index: number): void {
        setCharacterActions(
            [
                ...characterActions.slice(0, index + 1), 
                {character: "raiden", actions: ["normal attack"]},
                ...characterActions.slice(index + 1)
            ]
        )
    }
    
    function handleDeleteCharacterAction(index: number): void {
        setCharacterActions(characterActions.filter((_: object, idx: number) => idx !== index))
    }

    // TODO: refactor handleXyzAction functions to a single function. 
    // they only differ in what is being inserted at the specified index
    function handleAddAction(index: number): void {
        setCharacterActions(characterActions.map((ca: {character: string, actions: string[]}, idx) =>
            idx === index
             ? {character: ca.character, actions: [...ca.actions, "normal attack"]}
             : ca
            // at specified index, make a new characterAction (with an added action)   
            // otherwise map same old characterActions to new ones
            // equivalent to:
            // [
            //     ...characterActions.slice(0, index - 1),
            //     {   
            //         character: characterActions[index].character, 
            //         actions: [...characterActions[index].actions, "normal attack"]
            //     },
            //     ...characterActions.slice(index)
            // ]
        ))
    }
    function handleDeleteAction(index: number, actionIndex: number): void {
        setCharacterActions(characterActions.map((ca: {character: string, actions: string[]}, idx) =>
            idx === index
            ? {character: ca.character, actions: ca.actions.filter((_: string, aIdx: number) => aIdx !== actionIndex )}
            : ca
            // at specified index, make a new characterAction (with that action removed)   
            // otherwise map same old characterActions to new ones
        ))
    }
    function handleSelectAction(event: React.ChangeEvent<HTMLSelectElement>, index: number, actionIndex: number): void {
        console.log("handleSelectAction function called")
        const selectedAction = event.currentTarget.value
        setCharacterActions(characterActions.map((ca: {character: string, actions: string[]}, idx) =>
            idx === index
            ? {character: ca.character, actions: ca.actions.map((a: string, aIdx: number) => 
                aIdx === actionIndex ? selectedAction : a)}
            : ca
            // at specified index, make a new characterAction (with that action changed to the selected one)   
            // otherwise map same old characterActions to new ones
        )) 
    }
    function handleSelectCharacter(event: React.ChangeEvent<HTMLSelectElement>, index: number): void {
        console.log("triggered c")
        const selectedCharacter = event.target.value
        setCharacterActions(characterActions.map((ca: {character: string, actions: string[]}, idx) =>
            idx === index ? {character: selectedCharacter, actions: ca.actions} : ca
            // at specified index, make a new characterAction (with that character changed to the selected one)   
            // otherwise map same old characterActions to new ones
        ))
    }

    // TODO should be sent to ActionList cfg instead
    const rotation = characterActions.map(({character, actions}) => 
        <div>{character} {actions.toString().replace(",", ", ")}{";"}</div>
    )

    const actionsEl = (index: number, actions: string[]) =>
        <>
            {actions.map((action, actionIndex) =>    
                <div key={actionIndex}>
                    {/* delete action */}
                    <Button icon="cross" intent="danger" small onClick={() => handleDeleteAction(index, actionIndex)} />

                    {/* select action */}
                    <select
                        onChange={(e) => handleSelectAction(e, index, actionIndex)}
                        value={action}>
                        {availableActions.map((a, i) => <option key={i} value={a}>{a}</option>)} {/* TODO get actions from proper character or sth */}
                    </select>
                </div>
            )}

            {/* add action */}
            <button 
                style={{"borderRadius": "2px 2px"}}
                type="button"
                onClick={() => handleAddAction(index)}>
                <Icon icon="plus" size={20} color="white" />
            </button>
        </>

    const characterActionsEl = characterActions.map(({character, actions}, index) => 
        <div key={index} className="flex flex-row flex-wrap pl-2">
            <div style={{"border":"solid red 1px"}}>
                {/* delete CharacterAction */}
                <Button icon="cross" intent="danger" small onClick={() => handleDeleteCharacterAction(index)} />
                
                {/* select character */}
                <select 
                    onChange={e => handleSelectCharacter(e, index)}
                    value={characterActions[index].character}>                      
                    {availableCharacters.map((c, i) => <option key={i} value={c}>{c}</option>)} {/* TODO get characters from proper team */}
                </select>

                {actionsEl(index, actions)}
            </div>

            {/* add characterAction */}
            <button 
                style={{"borderRadius": "2px 2px"}}
                type="button"
                onClick={() => handleAddCharacterAction(index)}>
                <Icon icon="plus" size={20} color="white" />
            </button>
        </div>
    )

    return (
        <>
            {rotation}
            <div className="flex flex-row flex-wrap pl-2">
                {characterActionsEl}
            </div>
        </>
    )
}