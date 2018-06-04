import React from 'react';

export async function CreateRecord(record) {
    console.log("Create called")
    console.dir(record)
    console.log(JSON.stringify({
        "title": record.title,
        "location": {
            "lat": record.location.lat.toString(),   
            "lng": record.location.lng.toString(),
        }
    }))
    let response = await fetch('http://localhost:8000/v1/record', {
        method:'POST',
        headers: {'Content-Type':'application/json'},
    body: JSON.stringify({
        "title": record.title,
        "location": {
            "lat": record.location.lat,   
            "lng": record.location.lng,
        }
    })
    })

    return response
}