import { Outlet } from "react-router-dom"
import HeaderBar from "./HeaderBar"
import SidePanel from "./SidePanel"

export default function Layout(){
    return (
        <div>
            <HeaderBar/>
            <div style={{display:"flex", flexDirection:"row"}}>
                <SidePanel/>
                <div style={{flex:1,alignItems:"center"}}>
                <Outlet/>
                </div>
            </div>
            
        </div>
    )
}