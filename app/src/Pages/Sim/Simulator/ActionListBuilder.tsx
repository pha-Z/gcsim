import React from "react"
import { Button, Icon } from '@blueprintjs/core'
import { DragDropContext, Droppable, Draggable, DropResult, DraggableStateSnapshot } from "react-beautiful-dnd";
// dragdrop tutorial
// https://egghead.io/lessons/react-reorder-a-list-with-react-beautiful-dnd

// TODO get proper possible actions (which are also based on character?)
// likely needs to be changed entirely.
// Look into actions and characters from 'core', 'simActions' or whatever the thing is called
const availableActions: string[] = ["NA", "CA", "E", "Q", "D", "J"]
const availableCharacters: string[] = ["raiden", "xingqiu", "xiangling", "bennet"]
const defaultNewAction: string = "NA"
function defaultNewCharacterAction(): {character: string, actions: string[]} {
    return {character: availableCharacters[1], actions: [defaultNewAction]}
}



// TODO at hint descriptions for elements
export function ActionListBuilder() {
    const [characterActions, setCharacterActions] = React.useState(
        [
            {character: "raiden", actions: ["E", "D"]}, 
            defaultNewCharacterAction()
        ]
    )
    
    function handleAddCharacterAction(index: number): void {
        setCharacterActions(
            [
                ...characterActions.slice(0, index + 1), 
                defaultNewCharacterAction(),
                ...characterActions.slice(index + 1)
            ]
        )
    }
    
    function handleDeleteCharacterAction(index: number): void {
        setCharacterActions(characterActions.filter((_: object, idx: number) => idx !== index))
    }

    // TODO: refactor handleXyzAction functions to a single function (with a few helper functions). 
    // they only differ in what is being inserted at the specified index
    function handleAddAction(index: number): void {
        setCharacterActions(characterActions.map((ca: {character: string, actions: string[]}, idx) =>
            idx === index
             ? {character: ca.character, actions: [...ca.actions, defaultNewAction]}
             : ca
            // at specified index, replace with 'characterAction' (with an added 'action')   
            // otherwise map same old characterActions to new ones
            // equivalent to:
            // [
            //     ...characterActions.slice(0, index - 1),
            //     {   
            //         character: characterActions[index].character, 
            //         actions: [...characterActions[index].actions, defaultNewAction]
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
            // at specified index, replace with characterAction (with that action removed)   
            // otherwise map same old characterActions to new ones
        ))
    }
    function handleSelectAction(event: React.ChangeEvent<HTMLSelectElement>, index: number, actionIndex: number): void {
        const selectedAction = event.currentTarget.value
        setCharacterActions(characterActions.map((ca: {character: string, actions: string[]}, idx) =>
            idx === index
            ? {character: ca.character, actions: ca.actions.map((a: string, aIdx: number) => 
                aIdx === actionIndex ? selectedAction : a)}
            : ca
            // at specified index, replace with characterAction (with that action changed to the selected one)   
            // otherwise map same old characterActions to new ones
        )) 
    }
    function handleSelectCharacter(event: React.ChangeEvent<HTMLSelectElement>, index: number): void {
        const selectedCharacter = event.target.value
        setCharacterActions(characterActions.map((ca: {character: string, actions: string[]}, idx) =>
            idx === index ? {character: selectedCharacter, actions: ca.actions} : ca
            // at specified index, replace with characterAction (with that character changed to the selected one)   
            // otherwise map same old characterActions to new ones
        ))
    }


    // component building blocks
    const selectActionEl = (index: number, actionIndex: number, action: string) =>
        <select
            onChange={(e) => handleSelectAction(e, index, actionIndex)}
            value={action}
            style={{height: "40px", minHeight: "40px", borderRadius: "20px"}}
        >
            {availableActions.map((a, i) => 
                <option key={i} value={a}>{a}</option>)} {/* TODO get actions from proper character or sth */}
        </select>
    const deleteCharacterActionEl = (index: number) => 
        <div className="top-1 left-1">
            <Button 
                icon="cross"
                intent="danger"
                small
                onClick={() => handleDeleteCharacterAction(index)}
            />
        </div>
    const deleteActionEl = (index: number, actionIndex: number) => 
        <div className="top-1 left-1">
            <Button
                icon="cross" 
                intent="danger" 
                small 
                onClick={() => handleDeleteAction(index, actionIndex)}
            />
        </div>
    const selectCharacterEl = (index: number) =>
        <select 
            onChange={e => handleSelectCharacter(e, index)}
            value={characterActions[index].character}
            style={{height: "80px", minHeight: "80px", borderRadius: "40px"}}
        >                      
            {availableCharacters.map((c, i) =>  
                <option key={i} value={c}>{c}</option>)} {/* TODO get characters from proper team */}
        </select>
    const draggableActionsEl = (index: number, actions: string[]) => actions.map((action, actionIndex) => 
        <Draggable key={actionIndex} draggableId={index+"."+actionIndex} index={actionIndex}>
            {provided => 
                <div
                    ref={provided.innerRef} 
                    {...provided.draggableProps}
                    {...provided.dragHandleProps}
                    // cannot style this div or it bricks drag animations
                >
                    
                    
                    <div // styling div
                        style={{border: "solid red 1px", margin: "5px"}}
                    >
                        <div // styling div
                            style={{padding: "20px 5px", minWidth: "60px"}}
                            className="flex"
                        >
                            {deleteActionEl(index, actionIndex)}
                            {selectActionEl(index, actionIndex, action)}
                        </div>
                    </div>

                    
                </div>
            }
        </Draggable>     
    )
    const addActionEl = (index: number) =>
        <button 
            style={{borderRadius: "2px 2px"}}
            type="button"
            onClick={() => handleAddAction(index)}
        >
            <Icon icon="plus" size={20} color="white" />
        </button>
    const addCharacterActionEl = (index: number) =>
        <button 
            type="button"
            onClick={() => handleAddCharacterAction(index)}
        >
            <Icon icon="plus" size={30} color="white" />
        </button>
    const droppableActionsContainerEl = (index: number, actions: string[]) => 
        <Droppable droppableId={"actionsDropArea."+index} direction="horizontal" type={"dropAction"}>
            {provided =>
                <div 
                    ref={provided.innerRef}
                    {...provided.droppableProps}
                    style={{border: "solid lime 1px", minWidth: "30px"}}
                    className="flex"
                >

                    {deleteCharacterActionEl(index)}
                    {selectCharacterEl(index)}  
                    {draggableActionsEl(index, actions)}
                    {addActionEl(index)}

                    {provided.placeholder}
                </div>
            }
        </Droppable>
    const draggableCharacterActionsEl = characterActions.map((characterAction, index) => 
        <Draggable key={index} draggableId={index+""} index={index}>
            {provided => 
                <div 
                    ref={provided.innerRef}
                    {...provided.draggableProps}
                    {...provided.dragHandleProps}
                    // cannot stlye this div or it bricks drag animations
                >


                    <div className="flex"> 
                        <div // styling div
                            style={{
                                border: "solid green 1px",
                                borderRadius: "10px",
                                backgroundColor: "steelblue",
                                padding: "25px 5px",
                                margin: "5px 2px"
                            }}
                        >
                            {droppableActionsContainerEl(index, characterAction.actions)}
                        </div>
                        {addCharacterActionEl(index)}
                    </div>


                </div>
            }
        </Draggable>
    )
 
    
    // handle drag and drop functionality from beautiful-dnd
    function handleDrop(droppedItem: DropResult): void {
        if (!droppedItem.destination) {
            return
        }
        if (droppedItem.type === "dropCharacterAction") {
            const indexS = droppedItem.source.index
            const indexD = droppedItem.destination.index
    
            if (indexS === indexD) {
                return
            }
            
            const updatedCharacterActions = [...characterActions]
            const [reorderedItem] = updatedCharacterActions.splice(indexS, 1) // Remove dragged item
            updatedCharacterActions.splice(indexD, 0, reorderedItem) // Add dropped item
            
            setCharacterActions(updatedCharacterActions)
        }
        else if (droppedItem.type === "dropAction") {
            // droppableId={"actionsDropArea."+index}
            const IndexS = parseInt(droppedItem.source.droppableId.slice(16))
            const IndexD = parseInt(droppedItem.destination.droppableId.slice(16))
            const actionIndexS = droppedItem.source.index
            const actionIndexD = droppedItem.destination.index
    
            if (IndexS === IndexD && actionIndexS === actionIndexD) {
                return
            }
            
            const updatedCharacterActions = [...characterActions]

            const updatedActionsSource = updatedCharacterActions[IndexS].actions 
            const updatedActionsDestination = updatedCharacterActions[IndexD].actions
            // ^ note: both lists reference the same object
            // this is required
            // when the indices are the same, both reference the same 'actions'
            // this is redundant, but it doesnt matter 

            const [reorderedItem] = updatedActionsSource.splice(actionIndexS, 1) // Remove dragged item
            updatedActionsDestination.splice(actionIndexD, 0, reorderedItem) // Add dropped item
            
            setCharacterActions(characterActions.map((ca: {character: string, actions: string[]}, idx) =>
                idx === IndexS
                ? {character: ca.character, actions: updatedActionsSource}
                : idx === IndexD
                    ? {character: ca.character, actions: updatedActionsDestination}
                    : ca
                // at source index, replace with 'characterAction' (with the updated source 'actions')   
                // do the same thing at destination index  
                // otherwise map same old characterActions to new ones)
            ))
        }
    }



    // TODO should be sent to ActionList cfg instead
    const rotation = characterActions.map(({character, actions}) => 
        <div>{character} {actions.join(", ")}{";"}</div>
    )


    return (
        <>
            <div style={{minHeight: "200px", border: "solid white 1px"}}>
                {rotation}
            </div>
            {/* dragndrop wrapper */}
            <DragDropContext onDragEnd={handleDrop}>
                    <Droppable droppableId="characterActionsDropArea" direction="horizontal" type="dropCharacterAction">
                            {provided => 
                                <div
                                    ref={provided.innerRef}
                                    {...provided.droppableProps}
                                >

                                    <div // styling div
                                        className="flex"
                                        style={{border: "solid 5px grey", overflowX: "scroll"}}
                                    >
                                        {draggableCharacterActionsEl}
                                    </div>

                                    {provided.placeholder}
                                </div>
                            }
                    </Droppable>
            </DragDropContext>
        </>
    )
}