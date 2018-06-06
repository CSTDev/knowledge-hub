export function CreateRecord(record) {
    console.log("Create called")

    try{
        return fetch(process.env.REACT_APP_API_URL + '/record', {
            method:'POST',
            headers: {'Content-Type':'application/json'},
            body: JSON.stringify(
                record
            )
        }).then(response => {
            return response
        })
        
    } catch(err) {
        return null
    }
    
}