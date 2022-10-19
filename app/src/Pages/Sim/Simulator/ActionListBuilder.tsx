import React, { useEffect } from "react"
import { Button, Icon } from '@blueprintjs/core'
import { DragDropContext, Droppable, Draggable, DropResult, DraggableStateSnapshot } from "react-beautiful-dnd";
import { stringify } from "ajv";
import { defaultModifiers } from "@popperjs/core/lib/popper-lite";
// dragdrop tutorial
// https://egghead.io/lessons/react-reorder-a-list-with-react-beautiful-dnd

// TODO get proper possible actions (which are also based on character?)
// likely needs to be changed entirely.
// Look into actions and characters from 'core', 'simActions' or whatever the thing is called
const availableActions: string[] = ["NA", "NA:2", "NA:3", "NA:4", "NA:5", "CA", "E", "Q", "D", "J"]
const availableCharacters: string[] = ["raiden", "xingqiu", "xiangling", "bennett"]
const defaultNewAction: string = "NA"
function defaultNewCharacterAction(): {character: string, actions: string[]} {
    return {character: availableCharacters[1], actions: [defaultNewAction]}
}

// custom hook, to access previous state.
// https://stackoverflow.com/questions/53446020/how-to-compare-oldvalues-and-newvalues-on-react-hooks-useeffect
const usePrevious = <T extends unknown>(value: T): T | undefined => {
    const ref = React.useRef<T>()
    React.useEffect(() => {ref.current = value})
    return ref.current
}

// TODO at hint descriptions for elements
export function ActionListBuilder() {
    const [characterActions, setCharacterActions] = React.useState(
        [
            {character: "raiden", actions: ["E", "D"]}, 
            defaultNewCharacterAction()
        ]
    )
    // focus on select elements automatically while building the ActionList
    // --- currently focus is implemented with useEffect ---
    // --- could also implement focus with onKeypress event for 'enter' ---
    // keep reference to all select elements
    const refsToSelectElMap: Map<string, React.RefObject<HTMLSelectElement>> = new Map([])
    // after render, compare previous to new state it changed
    // and identify which select element should be focused next 
    const prevCharacterActions = usePrevious(characterActions)
    React.useEffect(() => {
        // do nothing
        // a. if there is no previous state (e.g. on page load)
        // b. if only a characterAction was deleted (characterActions.length must be smaller)
        // c. if its a drag and drop
        //    1 - the json length is unchaged in that case
        //    2 - but still more than 1 inner item changed
        // d. if the the state is still the same
        //    1 - the json length is unchanged in that case
        //    2 - and less than 1 inner item changed
        // e. c.1 is the same as d.1
        // f. c.2 && d2 is equivalent to "unequal 1 inner item changed"
        const isUnequal1Change = (): boolean => { 
            // note: type is always number ('false' is filtered out)
            const indicesOfDiffCharacterActions: number[] = characterActions.map((ca: {character: string, actions: string[]}, index: number) => 
                JSON.stringify(ca) !== JSON.stringify(prevCharacterActions[index]) && index
            ).filter(Number.isInteger)
            // return true if !== 1 characterAction changed
            if (indicesOfDiffCharacterActions.length !== 1) {
                return true
            }
            const indexOfDiff = indicesOfDiffCharacterActions[0]
            // return false if character changed
            if (characterActions[indexOfDiff].character !== prevCharacterActions[indexOfDiff].character) {
                return false
            }
            // note: type is always number ('false' is filtered out)
            const indicesOfDiffActions: number[] = characterActions[indexOfDiff].actions.map((action: string, actionIndex: number) => 
                action !== prevCharacterActions[indexOfDiff].actions[actionIndex] && actionIndex
            ).filter(Number.isInteger)
            // return true if !== 1 action changed
            if (indicesOfDiffActions.length !== 1) {
                return true
            }
            return false
        }
        if (
            !prevCharacterActions || // a.
            characterActions.length < prevCharacterActions.length || // b.
            (JSON.stringify(characterActions).length === JSON.stringify(prevCharacterActions).length && // e.
            (isUnequal1Change())) // f.
        ) { 
            return
        }
        // to identify the reference Key of the selectEL we want to focus
        let refKey: string = ""
        loop1: for (let index: number = 0; index < characterActions.length; index++) {
            if (
                // check for new selectCharEl at new index, which didnt exist previously
                // (need to check this first, because undefined.actions would throw an error)
                !prevCharacterActions[index] ||
                // check for new or changed selectCharacterEL 
                characterActions[index].character !== prevCharacterActions[index].character
            ) {
                // if its new, characterActions.length must be greater
                if (characterActions.length > prevCharacterActions.length) {
                    refKey = index.toString()
                    break
                } 
                // else its changed
                else {
                    // get key of the selectEl after that
                    // this is always the first (0th) selectActionEl at that index
                    // (doesnt matter when there arent any selectActionEl's)
                    refKey = index + "." + "0"
                    break
                }
            } 
            // else do nothing if only an action was deleted (actions.length must be smaller)
            else if (characterActions[index].actions.length < prevCharacterActions[index].actions.length) {
                return
            }
            // else check for new selectActionsEl (actions.length must be greater)
            else if (characterActions[index].actions.length > prevCharacterActions[index].actions.length) {
                // get key of the new selectActionEl
                // which is always the last action at that index
                const indexOfNewAction = characterActions[index].actions.length - 1
                refKey = index+"."+indexOfNewAction
                break
            }
            // else its a changed selectActionsEl
            else {
                for (
                    let actionIndex: number = 0; 
                    actionIndex < characterActions[index].actions.length; 
                    actionIndex++
                ) {
                    // position of changed selecActionEL
                    if (
                        characterActions[index].actions[actionIndex] !==
                        prevCharacterActions[index].actions[actionIndex]
                    ) {
                        // if its the last action at that index
                        if (actionIndex === characterActions[index].actions.length - 1) {
                            // get key of next selectCharacterEl at the next index
                            // (doesnt matter when there isnt one)
                            refKey = (index + 1).toString()
                            break loop1
                        }
                        // else get key of next selectActionsEl at that index
                        else {
                            refKey = index + "." + (actionIndex + 1)
                            break loop1
                        }
                    }
                }
            }
        }
        // focus select element
        refsToSelectElMap.get(refKey)?.current?.focus()
    }, [characterActions])


    // handlers
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
    const selectActionEl = (index: number, actionIndex: number, action: string) => {
        // create ref for every selectActionEl and store ref in refsToSelectElMap
        const refToSelectEl: React.RefObject<HTMLSelectElement> = React.createRef<HTMLSelectElement>()
        const key = index+"."+actionIndex
        refsToSelectElMap.set(key, refToSelectEl)

        return <select
            onChange={(e) => handleSelectAction(e, index, actionIndex)}
            value={action}
            style={{height: "52px", minHeight: "52px", borderRadius: "26px"}}
            ref={refToSelectEl}
        >
            {availableActions.map((a, i) => /* TODO get actions from proper character or sth */
                <option key={i} value={a}>{a}</option>
            )}
        </select>
    }
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
    const selectCharacterEl = (index: number) =>{
        const refToSelectEl: React.RefObject<HTMLSelectElement> = React.createRef<HTMLSelectElement>()
        const key = index.toString()
        refsToSelectElMap.set(key, refToSelectEl)
        
        return <select 
            onChange={e => handleSelectCharacter(e, index)}
            value={characterActions[index].character}
            style={{height: "80px", minHeight: "80px", borderRadius: "40px"}}
            ref={refToSelectEl}
        >             
            {availableCharacters.map((c, i) => /* TODO get characters from proper team */
                <option key={i} value={c}>{c}</option>
            )}
        </select>
    }
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
                        style={{border: "solid red 1px", margin: "10px 0"}}
                    >
                        <div // styling div
                            style={{padding: "21px 5px", minWidth: "50px"}}
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
            style={{padding: "0 20px"}}
            type="button"
            onClick={() => handleAddAction(index)}
        >
            <Icon icon="plus" size={20} color="white" />
        </button>
    const addCharacterActionEl = (index: number) =>
        <button 
            style={{marginLeft: "80px", padding: "1px 20px"}}
            type="button"
            onClick={() => handleAddCharacterAction(index)}
        >
            <Icon icon="plus" size={25} color="white" />
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
    const draggableCharacterActionsEl = characterActions.map((
        characterAction: {character: string, actions: string[]},
        index: number
    ) => 
        <div // styling div
            style={{marginLeft: (index * 30)+"px"}}
        >
            <Draggable draggableId={index+""} index={index}>
                {(provided, snapshot) => 
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
                                    padding: "21px 0",
                                    margin: "1px"
                                }}
                            >
                                {droppableActionsContainerEl(index, characterAction.actions)}
                            </div>
                        </div>
                        {addCharacterActionEl(index)}


                    </div>
                }
            </Draggable>
        </div>
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
    const rotation = 
        <div style={{minHeight: "200px", padding: "4px", border: "solid white 1px"}}>
            
            <div>while 1 {'{'}</div>
            
            {characterActions.map(({character, actions}, index) => 
                <div key={index}>{character} {actions.join(", ")}{";"}</div>
            )}
            
            <div>{'}'}</div>

        </div>


    return (
        <>
            {rotation}
            {/* dragndrop wrapper */}
            <DragDropContext onDragEnd={handleDrop}>
                    <Droppable droppableId="characterActionsDropArea" direction="vertical" type="dropCharacterAction">
                            {provided => 
                                <div
                                    ref={provided.innerRef}
                                    {...provided.droppableProps}
                                >

                                    <div // styling div
                                        style={{height: "1000px", border: "solid 5px grey", overflowX: "scroll", scrollBehavior: "smooth"}}
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