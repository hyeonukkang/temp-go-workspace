import {useState} from 'react';
import './App.css';


import {Greet, PostgresTestQuery} from "../wailsjs/go/main/App";
import { Button } from '@vuno/ui-antd';

function App() {
    const [resultText, setResultText] = useState(
      "Please enter your name below!!"
    )
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }



    return (
        <div id="App">
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <Button onClick={greet}>Greet</Button>
            </div>
            <div style={{marginTop: 24, textAlign: 'center'}}>
                <Button onClick={() => PostgresTestQuery().then(updateResultText)}>
                    PostgreSQL 연결 테스트
                </Button>
            </div>
        </div>
    )
}

export default App
