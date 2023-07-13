#!/usr/bin/env node
const arg= require("arg")


const inst=()=>{
    console.log('tool [CMD]\n --start\tstarts the cli\n --build\tbuils the app')
}

try{
    const args=arg({
        '--start' : Boolean,
        '--build' : Boolean
    })

    if (args['--start']){
        console.log("App is starting")
    }
}catch(error){
    console.log(error.message)
    console.log()
    inst()
}


