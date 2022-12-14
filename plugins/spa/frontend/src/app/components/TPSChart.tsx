import * as React from 'react';
import Card from "react-bootstrap/Card";
import NodeStore from "app/stores/NodeStore";
import {inject, observer} from "mobx-react";
import {Line} from "react-chartjs-2";
import {defaultChartOptions} from "app/misc/Chart";

interface Props {
    nodeStore?: NodeStore;
}

const lineChartOptions = Object.assign({
    scales: {
        xAxes: [{
            ticks: {
                autoSkip: true,
                maxTicksLimit: 8,
                fontSize: 8,
                minRotation: 0,
                maxRotation: 0,
            },
            gridLines: {
                display: false
            }
        }],
        yAxes: [{
            gridLines: {
                display: false
            },
            ticks: {
                fontSize: 10,
                maxTicksLimit: 4,
                beginAtZero: true,
            },
        }],
    },
}, defaultChartOptions);

@inject("nodeStore")
@observer
export default class TPSChart extends React.Component<Props, any> {
    render() {
        return (
            <Card>
                <Card.Body>
                    <Card.Title>Transactions Per Second</Card.Title>
                    <small>
                        Incoming: {this.props.nodeStore.last_tps_metric.incoming}.
                        Outgoing: {this.props.nodeStore.last_tps_metric.outgoing}.
                        New: {this.props.nodeStore.last_tps_metric.new}.
                    </small>

                    <Line height={50} data={this.props.nodeStore.tpsSeries} options={lineChartOptions}/>
                </Card.Body>
            </Card>
        );
    }
}
