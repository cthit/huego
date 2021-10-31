import "../App.css"
import React from "react"
import {toHslString} from "../utility";

export const LightStrip = props => {
    let hsb = {
        h: 0,
        s: 100,
        b: 100
    }
    if (props.color.hsb) {
        hsb = props.color.hsb
    }

    return (
        <div className="LightStrip" style={{
            backgroundColor: toHslString(hsb)
        }}/>
    )
}