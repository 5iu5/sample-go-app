import {useNavigate} from 'react-router-dom'

export default function HeaderBar(){
    const navigate = useNavigate()

    return(
        <div style={{width:"100%"}}>
            <div style={{display: "flex", flexDirection: "row", justifyContent:"space-between", width:"100%", paddingLeft: "20px", paddingRight: "20px"
            }}>
                <h3 onClick={()=>{
                    console.log("home clicked")
                    navigate(`/`)
                    }}
                    style={{cursor:"pointer"}}>
                        logo</h3>
                <h3>search bar</h3>
                <h3>settings</h3>
            </div>
            <div style={{width:"100%", height:"1px", backgroundColor:"grey"}}/>
        </div>

    )
}


