import { useTheme } from '@emotion/react';
import { Box, Button, Drawer, Grid, Typography } from '@mui/material';
import { useSnackbar } from 'notistack';
import { useEffect, useRef, useState } from 'react';
import ReactFlow, { addEdge, Controls } from 'react-flow-renderer';
import { useHistory, useParams, useLocation } from 'react-router-dom';
import CustomChip from '../../components/CustomChip';
import CustomLine from '../../components/CustomNodesContent/CustomLine';
import PublishPipelineDrawer from '../../components/DrawerContent/PublishPipelineDrawer';
import RemoveLogsPageItem from '../../components/MoreInfoContent/RemoveLogsPageItem';
import MoreInfoMenu from '../../components/MoreInfoMenu';
import { useGetPipelineFlow } from '../../graphql/getPipelineFlow';
import { edgeTypes, nodeTypes, useGlobalFlowState } from '../Flow';
import { useGlobalEnvironmentState } from '../../components/EnviromentDropdown';
import RunsDropdown from './RunsDropdown';
import Timer from './Timer';

const View = () => {
    // Hooks
    const theme = useTheme();
    const history = useHistory();
    const { state: pipeline } = useLocation();
    const getPipelineFlow = useGetPipelineFlowHook(pipeline);

    // URI parameter
    const { pipelineId } = useParams();

    // Global states
    const FlowState = useGlobalFlowState();

    // Page states
    const [isOpenPublishDrawer, setIsOpenPublishDrawer] = useState(false);
    const [, setIsLoadingFlow] = useState(true);
    // // const [timeElapsed, setTimeElapsed] = useState(Math.floor((new Date() - FlowState.startedRunningAt.attach(Downgraded).get()) / 1000));
    // const [timeElapsed, setTimeElapsed] = useState(0);
    // const [tasks, setTasks] = useState();

    const Environment = useGlobalEnvironmentState();

    //Offset states and refs
    const [offsetHeight, setOffsetHeight] = useState(0);
    const offsetRef = useRef(null);

    useEffect(() => {
        setOffsetHeight(offsetRef.current.clientHeight);
    }, [offsetRef]);

    // Flow states
    const reactFlowWrapper = useRef(null);
    const [, setReactFlowInstance] = useState(null);
    const [elements, setElements] = useState([]);

    // Fetch previous elements
    useEffect(() => {
        // const prevElements = FlowState.elements.get();
        // setElements([...prevElements]);

        getPipelineFlow(Environment.id.get(), setElements);
        setIsLoadingFlow(false);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // Handle edit button
    const handleGoToEditorPage = () => {
        FlowState.isEditorPage.set(true);
        history.push({ pathname: `/pipelines/flow/${pipelineId}`, state: pipeline });
    };

    //Flow methods
    const onLoad = (_reactFlowInstance) => setReactFlowInstance(_reactFlowInstance);
    const onConnect = (params) => {
        setElements((els) => addEdge({ ...params, type: 'custom' }, els));
    };

    // // Updates timer every second
    // useEffect(() => {
    //     let secTimer;
    //     if (FlowState.isRunning.get()) {
    //         secTimer = setInterval(() => {
    //             setTimeElapsed(Math.floor((new Date() - FlowState.startedRunningAt.attach(Downgraded).get()) / 1000));
    //         }, 1000);
    //     }

    //     return () => {
    //         clearInterval(secTimer);
    //         setTimeElapsed(0);
    //     };
    // }, [FlowState.isRunning.get()]);

    return (
        <Box className="page" height="calc(100vh - 100px)" minHeight="min-content">
            <Box ref={offsetRef}>
                <Grid container alignItems="center" justifyContent="space-between" wrap="nowrap">
                    <Box display="flex">
                        <Typography component="h2" variant="h2" color="text.primary">
                            Pipelines {'>'} {pipeline.name}
                        </Typography>

                        <Grid display="flex" alignItems="flex-start">
                            <Box display="flex" alignItems="center" ml={4} mr={4}>
                                <Box height={16} width={16} backgroundColor={pipeline.online ? 'status.pipelineOnlineText' : 'error.main'} borderRadius="100%"></Box>
                                <Typography ml={1} fontSize={16} color={pipeline.online ? 'status.pipelineOnlineText' : 'error.main'}>
                                    {pipeline.online ? 'Online' : 'Offline'}
                                </Typography>
                            </Box>

                            <Box mr={3} textAlign="right">
                                <Typography variant="h3">{LOGS_MOCK.last_run}</Typography>
                                <Typography fontSize={17}>Last run</Typography>
                            </Box>

                            <Box ml={3} mr={4.5} textAlign="right">
                                <Typography variant="h3">{LOGS_MOCK.version}</Typography>
                                <Typography fontSize={17}>Version</Typography>
                            </Box>

                            <Button variant="contained" onClick={handleGoToEditorPage}>
                                Edit
                            </Button>
                            <Button variant="contained" sx={{ ml: 2.4 }} onClick={() => setIsOpenPublishDrawer(true)}>
                                Publish
                            </Button>
                            <Box sx={{ position: { xxs: 'relative', xl: 'absolute' }, ml: { xxs: 2, xl: 0 }, top: '0', right: '0' }}>
                                <MoreInfoMenu>
                                    <RemoveLogsPageItem />
                                </MoreInfoMenu>
                            </Box>
                        </Grid>
                    </Box>
                </Grid>

                <Grid mt={4} container alignItems="center" sx={{ width: { xl: '88%' }, flexWrap: 'nowrap' }}>
                    <Grid item display="flex" alignItems="center" sx={{ alignSelf: 'center' }} flex={1.2}>
                        <CustomChip amount={2} label="Steps" margin={2} customColor="orange" />
                        <CustomChip amount={2} label="Running" margin={2} customColor="blue" />
                        <CustomChip amount={2} label="Succeeded" margin={2} customColor="green" />
                        <CustomChip amount={2} label="Failed" margin={2} customColor="red" />
                        <CustomChip amount={2} label="Workers online" margin={2} customColor="purple" />
                    </Grid>

                    {/* Runs dropdown */}
                    <RunsDropdown />

                    {/* Timer */}
                    <Timer environmentID={Environment.id.get()} />
                </Grid>
            </Box>

            <Box mt={7} sx={{ position: 'absolute', top: offsetHeight, left: 0, right: 0, bottom: 0 }} ref={reactFlowWrapper}>
                {elements && elements.length > 0 ? (
                    <ReactFlow
                        nodeTypes={nodeTypes}
                        elements={elements}
                        onLoad={onLoad}
                        onConnect={onConnect}
                        connectionLineComponent={CustomLine}
                        edgeTypes={edgeTypes}
                        arrowHeadColor={theme.palette.mode === 'dark' ? '#fff' : '#222'}
                        snapToGrid={true}
                        snapGrid={[15, 15]}>
                        <Controls style={{ left: 'auto', right: 10 }} />
                    </ReactFlow>
                ) : (
                    <Box sx={{ height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                        <Typography>Create a pipeline by dragging the components here</Typography>
                    </Box>
                )}
            </Box>

            <Drawer anchor="right" open={isOpenPublishDrawer} onClose={() => setIsOpenPublishDrawer(!isOpenPublishDrawer)}>
                <PublishPipelineDrawer handleClose={() => setIsOpenPublishDrawer(false)} />
            </Drawer>
        </Box>
    );
};

const LOGS_MOCK = {
    id: 1,
    online: true,
    last_run: '22 Nov 2021 08:00',
    version: '0.0.1',
};

export default View;

// ------ Custom hooks
const useGetPipelineFlowHook = (pipeline) => {
    // GraphQL hook
    const getPipelineFlow = useGetPipelineFlow();

    // React router
    const history = useHistory();

    // Global state
    const FlowState = useGlobalFlowState();

    // URI parameter
    const { pipelineId } = useParams();

    const { enqueueSnackbar, closeSnackbar } = useSnackbar();

    // Get members
    return async (environmentID, setElements) => {
        const rawResponse = await getPipelineFlow({ pipelineID: pipelineId, environmentID });
        const response = prepareInputForFrontend(rawResponse);

        if (response.length === 0) {
            FlowState.elements.set([]);
            history.push({ pathname: `/pipelines/flow/${pipelineId}`, state: pipeline });
        } else if (response.r === 'error') {
            closeSnackbar();
            enqueueSnackbar("Can't get flow: " + response.msg, { variant: 'error' });
        } else if (response.errors) {
            response.errors.map((err) => enqueueSnackbar(err.message + ': get flow', { variant: 'error' }));
        } else {
            setElements(response);
            FlowState.elements.set(response);
        }
    };
};

// ----- Utility functions
function prepareInputForFrontend(input) {
    const edgesInput = [];
    const nodesInput = [];

    for (const edge of input.edges) {
        edgesInput.push({
            source: edge.from,
            sourceHandle: edge.meta.sourceHandle,
            target: edge.to,
            targetHandle: edge.meta.targetHandle,
            type: edge.meta.edgeType,
            arrowHeadType: edge.meta.arrowHeadType,
            id: edge.edgeID,
        });
    }

    for (const node of input.nodes) {
        let data = {
            ...node.meta?.data,
            name: node.name,
            description: node.description,
            workerGroup: node.workerGroup,
            commands: node.commands,
        };
        nodesInput.push({
            id: node.nodeID,
            type: node.nodeTypeDesc + 'Node',
            position: {
                x: node.meta.position.x,
                y: node.meta.position.y,
            },
            data,
        });
    }

    return [...edgesInput, ...nodesInput];
}
