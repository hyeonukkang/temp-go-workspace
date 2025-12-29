import {useEffect, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {GetFileInfoTable, LogButtonClick} from "../wailsjs/go/main/App";
import { Button } from '@vuno/ui-antd';

function App() {
    const [rows, setRows] = useState<Array<{id: number, file_name: string, uploaded_at: string}>>([]);
    const fetchTable = () => {
        GetFileInfoTable().then((data: any) => {
            if (Array.isArray(data)) {
                setRows(data.map((row) => ({
                    id: row.id,
                    file_name: row.file_name,
                    uploaded_at: row.uploaded_at
                })));
            }
        });
    };
    useEffect(() => {
        fetchTable();
        const interval = setInterval(fetchTable, 2000); // 2초마다 polling
        return () => clearInterval(interval);
    }, []);

    const handleClick = () => { 
        LogButtonClick(); // go 서버에 로그 기록
    }

    return (
        <div id="App">
                <Button onClick={handleClick}>file_info 테이블 목록</Button>
            <button className="btn" style={{marginBottom: 16}} onClick={fetchTable}>새로고침</button>
            <table style={{margin: '0 auto', minWidth: 400, borderCollapse: 'collapse'}}>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>파일명</th>
                        <th>업로드 시간</th>
                    </tr>
                </thead>
                <tbody>
                    {rows.map(row => (
                        <tr key={row.id}>
                            <td>{row.id}</td>
                            <td>{row.file_name}</td>
                            <td>{row.uploaded_at}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

export default App
