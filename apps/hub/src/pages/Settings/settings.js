import React, { Component } from 'react'
import { ToastContainer, toast } from 'react-toastify';
import { Paper, Button, TextField } from 'react-md';
import ReactDOM from 'react-dom';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';

import MenuBar from '../../components/MenuBar/menuBar';
import { VersionBar } from '../../components/VersionBar/versionBar';
import { UpdateFields, UpdateField, LoadFields } from '../../data/api';

import './settings.css'

const getItems = count =>
    Array.from({ length: count }, (v, k) => k).map(k => ({
        id: `item-${k}`,
        content: `item ${k}`,
    }));

const reorder = (list, startIndex, endIndex) => {
    const result = Array.from(list);
    const [removed] = result.splice(startIndex, 1);
    result.splice(endIndex, 0, removed);

    return result;
};

class Settings extends Component {

    constructor(props) {
        super(props);

        this.state = {
            version: process.env.REACT_APP_VERSION ? process.env.REACT_APP_VERSION : "0.0.1",
            items: []
        };

        this.onDragEnd = this.onDragEnd.bind(this)
    }

    componentDidMount = () => {
        LoadFields().then((response) => {                

            if(!response || response.status !== 200){
                console.dir(response)
                console.log(response.message)
                if (response !== null && response.message == 404){
                    toast("No fields set")
                    return
                }
              toast("Failed to load fields")
              return
            }
            response.json().then(json => {
                return json
              }).then(fieldList =>{
                this.setState({items: fieldList})
              })
            
        });
    }

    onDragEnd(result) {
        // dropped outside the list
        if (!result.destination) {
            return;
        }

        const items = reorder(
            this.state.items,
            result.source.index,
            result.destination.index
        );

        for(let i = 0; i < items.length; i++){
            items[i].order = i
        }

        this.setState({
            items,
        });
    }

    addItem = () => {
        let items = this.state.items;
        const item = {
            id: `item-` + items.length, //TODO Change this to be more unique
            value: `item ` + items.length, //TODO Just be empty on creation
            order: items.length + 1,
        }
        items.push(item);
        this.setState({
            items,
        });
    }

    deleteItem = (id) => {
        let itemsArray = [...this.state.items];
        
        itemsArray.splice(id - 1, 1);
        this.setState({
            items: itemsArray
        });
    }

    saveChanges = () => {
        UpdateFields(this.state.items);
    }

    onFieldUpdate = (itemId, e) => {
        console.log(itemId);
        const value = e.target.value;
        UpdateField(itemId, value);
    }

    render() {

        return (
            <div>
                <ToastContainer />
                <VersionBar version={this.state.version} />
                <MenuBar homeAction={() => this.saveChanges()}/>
                <Paper className="settingsContent" zDepth={3}>
                    <h1>Settings</h1>
                    <DragDropContext onDragEnd={this.onDragEnd} className="fieldList">
                        <Droppable droppableId="droppable">
                            {(provided, snapshot) => (
                                <div
                                    ref={provided.innerRef}

                                >
                                    {this.state.items.map((item, index) => (
                                        <Draggable key={item.id} draggableId={item.id} index={index} >
                                            {(provided, snapshot) => (
                                                <div
                                                    ref={provided.innerRef}
                                                    {...provided.draggableProps}
                                                    {...provided.dragHandleProps}
                                                    className="fieldItem"
                                                >
                                                    <TextField
                                                        className="settingsInput"
                                                        id={item.id}
                                                        defaultValue={item.value}
                                                        onBlur={(e) => this.onFieldUpdate(item.id,e)}
                                                    ></TextField>
                                                    <Button className="settings__delete-field" icon onClick={() => this.deleteItem(item.id)}>close</Button>
                                                </div>
                                            )}
                                        </Draggable>
                                    ))}
                                    {provided.placeholder}
                                </div>
                            )}
                        </Droppable>
                    </DragDropContext>
                    <div className="settings__button-div">
                        <Button raised primary className="settings__save-button" onClick={() => this.saveChanges()}>Save</Button>
                        <Button floating primary className="settings__add-button" onClick={() => this.addItem()}>add</Button>
                    </div>
                </Paper>
            </div>
        )
    }
}

export default Settings