import React from "react"

// TODO get proper possible actions (which is also based on character?)
// likely needs to be changed.
// Look into actions and characters from 'core', 'simActions' or whatever the thing is called
const possibleActionList: string[] = [
    "normal attack",
    "charge attack",
    "skill",
    "burst",
    "dash",
    "jump"
]

// An action performed by a character in the team
// which is to be added to the ActionList cfg
export function Action(
    props:{characters: string[], action?: string}  // TODO make sure we pass the proper characters/team list
) {
    const [character, setCharacter] = React.useState(props.characters[0]) // TODO make sure we pass the proper characters/team list
    const [action, setAction] = React.useState(props.action || possibleActionList[0])

    function switchCharacter(event: React.ChangeEvent<HTMLSelectElement>): void {
        const selectedCharacter = event.target.value
        setCharacter(selectedCharacter)
    }

    function switchAction(event: React.ChangeEvent<HTMLSelectElement>): void {
        const selectedAction = event.target.value
        setAction(selectedAction)
    }

    const possibleActionOptions = possibleActionList.map(a => (<option value={a}>{a}</option>)) 

    return (
        <div style={{"border":"solid red 1px"}}>
            <select
            onSelect={switchCharacter}>
                {/* TODO: loop through team members to give real options*/}
                <option value={character}>{character}</option>
            </select>
            <select
            onSelect={switchAction}>
                {possibleActionOptions}
            </select>
            <br/>
        </div>
    )
    }