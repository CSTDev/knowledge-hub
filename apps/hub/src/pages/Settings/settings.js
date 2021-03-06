import React, { Component } from 'react'
import { ToastContainer, toast } from 'react-toastify';
import { Paper, Button, TextField } from 'react-md';
import ReactDOM from 'react-dom';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import { v4 as uuid } from 'uuid';

import MenuBar from '../../components/MenuBar/menuBar';
import { VersionBar } from '../../components/VersionBar/versionBar';
import { UpdateFields, DeleteField, LoadFields } from '../../data/api';

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
            version: window.APP_CONFIG.VERSION ? window.APP_CONFIG.VERSION : "0.0.1",
            items: [],
            disableEdit: false
        };

        this.onDragEnd = this.onDragEnd.bind(this)
    }

    componentDidMount = () => {
        LoadFields().then((response) => {

            if (!response || response.status !== 200) {
                if (response !== null && response.message == 404) {
                    toast("No fields set")
                    return
                }
                toast("Failed to load fields")
                this.setState({ disableEdit: true })
                return
            }
            response.json().then(json => {
                return json
            }).then(fieldList => {
                let sortedfields = [].concat(fieldList).sort((a, b) => a.order > b.order)
                this.setState({ items: sortedfields })
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

        for (let i = 0; i < items.length; i++) {
            items[i].order = i
        }

        this.setState({
            items,
        });
    }

    addItem = () => {
        let items = this.state.items;
        const item = {
            id: uuid(),
            value: "",
            order: items.length,
        }
        items.push(item);
        this.setState({
            items,
        });
    }

    deleteItem = (id) => {
        let itemsArray = [...this.state.items];
        DeleteField(itemsArray[id].id).then(response => {
            if (!response || (response.status !== 200 && response.message !== "404")) {
                toast("Unable to remove field")
                return
            } else {
                itemsArray.splice(id, 1);

                for (let i = 0; i < itemsArray.length; i++) {
                    itemsArray[i].order = i
                }

                this.setState({
                    items: itemsArray
                });
            }
        });


    }

    saveChanges = async () => {
        if (this.validFields()) {
            let response = await UpdateFields(this.state.items);
            if (!response || response.status !== 200) {
                toast("Failed to save changes")
                return false
            } else {
                toast("Fields saved")
                return true
            }

        }
    }

    validFields = () => {
        let allValid = true
        this.state.items.map(field => {
            if (field.value == "") {
                toast("Field value cannot be empty");
                allValid = false;
            }
        })
        return allValid
    }

    onFieldUpdate = (itemId, e) => {
        let itemsArray = [...this.state.items];
        var newFields = itemsArray.map(field => {
            if (field.id == itemId) {
                return Object.assign({}, field, { value: e.target.value })
            }
            return field
        });
        this.setState({ items: newFields })
    }

    render() {

        return (
            <div>
                <ToastContainer />
                <VersionBar version={this.state.version} />
                <MenuBar homeAction={() => this.saveChanges()} />
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
                                                        placeholder="Enter field name"
                                                        onBlur={(e) => this.onFieldUpdate(item.id, e)}
                                                    ></TextField>
                                                    <Button className="settings__delete-field" icon onClick={() => { if (window.confirm('Are you sure you wish to delete this field? No one will be able to see the data associated with it.')) this.deleteItem(item.order) }}>close</Button>
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
                        <Button raised primary className="settings__save-button" disabled={this.state.disableEdit} onClick={() => this.saveChanges()}>Save</Button>
                        <Button floating primary className="settings__add-button" disabled={this.state.disableEdit} onClick={() => this.addItem()}>add</Button>
                    </div>
                </Paper>
            </div>
        )
    }
}

export default Settings