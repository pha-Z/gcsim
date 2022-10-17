import React from "react"
import { Action } from "./Action"

export function ActionListBuilder() {
    // TODO pass proper characters/team to <Action>
    const [actions, setActions] = React.useState([<Action key={0} characters={["raiden"]} action={"skilll"}/>])

    function addNewAction(): void {
        // TODO pass proper characters/team to <Action>
        setActions(currentActions => [...currentActions, <Action key={currentActions.length} characters={["raiden"]}/> ])
        console.log(actions)
    }   

    return (
        <>
            <div className="flex flex-row flex-wrap pl-2">
                {/* a list of <Action>'s, that each add a line to the cfg
                this component needs a button to append or delete an <Action> from the list */}
                {actions}
                <button 
                type="button"
                onClick={addNewAction}>+</button>
            </div>
        </>
    )
}