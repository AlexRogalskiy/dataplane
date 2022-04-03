import { ActionLayer } from './ActionLayer';
import { useTheme } from '@emotion/react';
import { faExpandArrowsAlt } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Box, Chip, Drawer, Grid, Typography } from '@mui/material';
import { useSnackbar } from 'notistack';
import { useEffect, useRef, useState } from 'react';
import ReactFlow, { addEdge, ControlButton, Controls, ReactFlowProvider } from 'react-flow-renderer';
import { useParams } from 'react-router-dom';
import CustomLine from '../../../components/CustomNodesContent/CustomLine';
import PublishPipelineDrawer from '../../../components/DrawerContent/PublishPipelineDrawer';
import { useGlobalEnvironmentState } from '../../../components/EnviromentDropdown';
import MenuItem from './MenuItem';
import MoreInfoMenu from '../../../components/MoreInfoMenu';
import { useGetPipelineFlow } from '../../../graphql/getPipelineFlow';
import { edgeTypes, nodeTypes, useGlobalFlowState } from '../../Flow';
import LogsDrawer from '../../../components/DrawerContent/LogsDrawer';
import TurnOffPipelineDrawer from '../../../components/DrawerContent/TurnOffPipelineDrawer';
import CustomChip from '../../../components/CustomChip';
import { Analytics } from './Analytics';
import { useGetActiveDeployment } from '../../../graphql/getActiveDeployment';

const DeploymentView = () => {
    const Environment = useGlobalEnvironmentState();

    // Hooks
    const theme = useTheme();
    const [deployment, setDeployment] = useState(null);
    const getPipelineFlow = useGetPipelineFlowHook();
    const getActiveDeployment = useGetActiveDeploymentHook(Environment.id.get(), setDeployment);

    // Global states
    const FlowState = useGlobalFlowState();

    // URI parameter
    const { version } = useParams();

    // Page states
    const [isOpenPublishDrawer, setIsOpenPublishDrawer] = useState(false);
    const [, setIsLoadingFlow] = useState(true);
    const [isOpenAnalytics, setIsOpenAnalytics] = useState(false);

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
    const [panOnDrag, setPanOnDrag] = useState(FlowState.isPanEnable.get());

    useEffect(() => {
        setPanOnDrag(FlowState.isPanEnable.get());

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [FlowState.isPanEnable.get()]);

    // Fetch previous elements
    useEffect(() => {
        if (!Environment.id.get()) return;
        setIsLoadingFlow(false);
        getActiveDeployment();

        document.querySelector('#root div').scrollTo(0, 0);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [Environment.id.get()]);

    useEffect(() => {
        getPipelineFlow(Environment.id.get(), setElements);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [FlowState.isEditorPage.get()]);

    // Trigger the scale button on keyboard 's' key click
    useEffect(() => {
        window.addEventListener('keydown', handleKeyDown);
        return () => {
            window.removeEventListener('keydown', handleKeyDown);
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const handleKeyDown = (e) => {
        if (e.keyCode === 83) {
            onZoomActive();
        }
    };

    //Flow methods
    const onLoad = (_reactFlowInstance) => setReactFlowInstance(_reactFlowInstance);
    const onConnect = (params) => {
        setElements((els) => addEdge({ ...params, type: 'custom' }, els));
    };
    const onConnectStart = () => {
        FlowState.isDragging.set(true);
        document.body.style.cursor = 'grabbing';
    };
    const onConnectEnd = () => {
        FlowState.isDragging.set(false);
        document.body.style.cursor = 'default';
    };
    const onMoveStart = (flow) => {
        FlowState.scale.set(flow.zoom);
    };
    const onMoveEnd = (flow) => {
        FlowState.scale.set(flow.zoom);
    };
    const onZoomActive = () => {
        FlowState.isPanEnable.set(!panOnDrag);

        if (panOnDrag) {
            document.body.style.cursor = 'default';
        } else {
            document.body.style.cursor = 'move';
        }
    };

    return (
        <Box className="page" height="calc(100vh - 136px)" minHeight="min-content">
            <Box ref={offsetRef}>
                <Grid container alignItems="center" justifyContent="space-between" wrap="nowrap">
                    <Box display="flex" alignItems="center">
                        <Typography component="h2" variant="h2" color="text.primary">
                            Deployments {'>'} {deployment?.name}
                        </Typography>

                        <Grid display="flex" alignItems="center">
                            <Box display="flex" alignItems="center" ml={4} mr={2}>
                                {deployment?.online ? <CustomChip label={'Online'} customColor="green" /> : <CustomChip label="Offline" customColor="red" />}
                            </Box>

                            <Chip
                                style={{
                                    borderRadius: 5,
                                    marginLeft: 5,
                                    marginRight: 10,
                                    fontWeight: 700,
                                    backgroundColor: deployment?.deploy_active && deployment.version === version ? '#7B61FF' : '#B9B9B9',
                                    color: '#FFF',
                                }}
                                label={`Deployed v${version ? version : deployment?.version}`}
                            />

                            {deployment?.version !== version ? (
                                <Typography ml={2} mr={2} variant="body2" color="editorPage.fileManagerIcon">
                                    Previous version
                                </Typography>
                            ) : null}

                            <Box sx={{ top: '0', right: '0' }}>
                                <MoreInfoMenu iconHorizontal>
                                    <MenuItem
                                        pipeline={deployment}
                                        getPipelineFlow={() => getPipelineFlow(Environment.id.get(), setElements)}
                                        isPipelineOnline={deployment?.online}
                                        getPipeline={getActiveDeployment}
                                        setIsOpenAnalytics={setIsOpenAnalytics}
                                    />
                                </MoreInfoMenu>
                            </Box>
                        </Grid>
                    </Box>
                </Grid>

                {/* Run/Stop button, Chips, Timer */}
                <ActionLayer setElements={setElements} environmentId={Environment.id.get()} deployment={deployment} />
            </Box>
            {!FlowState.isOpenLogDrawer.get() && !isOpenAnalytics ? (
                <Box mt={7} sx={{ position: 'absolute', top: offsetHeight, left: 0, right: 0, bottom: 0 }} ref={reactFlowWrapper}>
                    {elements && elements.length > 0 ? (
                        <ReactFlowProvider>
                            <ReactFlow
                                zoomOnScroll={false}
                                zoomOnPinch={false}
                                paneMoveable={panOnDrag || false}
                                onMoveStart={onMoveStart}
                                onMoveEnd={onMoveEnd}
                                nodeTypes={nodeTypes}
                                elements={elements}
                                defaultZoom={FlowState.scale.get()}
                                nodesDraggable={false}
                                nodesConnectable={false}
                                preventScrolling={false}
                                onLoad={onLoad}
                                onConnect={onConnect}
                                onConnectStart={onConnectStart}
                                onConnectEnd={onConnectEnd}
                                connectionLineComponent={CustomLine}
                                edgeTypes={edgeTypes}
                                arrowHeadColor={theme.palette.mode === 'dark' ? '#fff' : '#222'}
                                snapToGrid={true}
                                snapGrid={[15, 15]}>
                                <Controls style={{ left: 'auto', right: 10, bottom: 50 }}>
                                    <ControlButton onClick={onZoomActive} style={{ border: `1px solid ${FlowState.isPanEnable.get() ? '#72B842' : 'transparent'}` }}>
                                        <Box component={FontAwesomeIcon} icon={faExpandArrowsAlt} sx={{ color: FlowState.isPanEnable.get() ? '#72B842' : '' }} />
                                    </ControlButton>
                                </Controls>
                                <Box sx={{ position: 'absolute', left: 'auto', right: 10, bottom: 10 }}>
                                    <Typography fontSize={12}>Scale {Math.floor((FlowState.scale.get() || 1) * 100)}%</Typography>
                                </Box>
                            </ReactFlow>
                        </ReactFlowProvider>
                    ) : (
                        <Box sx={{ height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                            <Typography>Create a pipeline by dragging the components here</Typography>
                        </Box>
                    )}
                </Box>
            ) : null}

            {isOpenAnalytics ? <Analytics setIsOpenAnalytics={setIsOpenAnalytics} /> : null}
            <Drawer anchor="right" open={isOpenPublishDrawer} onClose={() => setIsOpenPublishDrawer(!isOpenPublishDrawer)}>
                <PublishPipelineDrawer handleClose={() => setIsOpenPublishDrawer(false)} />
            </Drawer>
            <Drawer
                hideBackdrop
                sx={{
                    width: 'calc(100% - 203px)',
                    height: 'calc(100% - 82px)',
                    [`& .MuiDrawer-paper`]: { width: 'calc(100% - 203px)', top: 82, height: 'calc(100% - 82px)', background: '#222', paddingBottom: 2 },
                }}
                anchor="right"
                open={FlowState.isOpenLogDrawer.get()}
                onClose={() => FlowState.isOpenLogDrawer.set(false)}>
                <LogsDrawer handleClose={() => FlowState.isOpenLogDrawer.set(false)} environmentId={Environment.id.get()} />
            </Drawer>
            <Drawer anchor="right" open={FlowState.isOpenTurnOffPipelineDrawer.get()} onClose={() => FlowState.isOpenTurnOffPipelineDrawer.set(false)}>
                <TurnOffPipelineDrawer
                    handleClose={() => FlowState.isOpenTurnOffPipelineDrawer.set(false)} //
                    pipelineID={deployment?.pipelineID}
                    environmentID={deployment?.environmentID}
                    name={deployment?.name}
                    getPipelineFlow={() => getPipelineFlow(Environment.id.get(), setElements)}
                    getPipeline={getActiveDeployment}
                />
            </Drawer>
        </Box>
    );
};

export default DeploymentView;

// ------ Custom hooks
export const useGetPipelineFlowHook = () => {
    // GraphQL hook
    const getPipelineFlow = useGetPipelineFlow();

    // Global state
    const FlowState = useGlobalFlowState();

    // URI parameter
    const { deploymentId } = useParams();

    const pipelineID = deploymentId.replace('d-', '');

    const { enqueueSnackbar, closeSnackbar } = useSnackbar();

    // Get members
    return async (environmentID, setElements) => {
        const rawResponse = await getPipelineFlow({ pipelineID, environmentID });
        const response = prepareInputForFrontend(rawResponse);

        if (response.length === 0) {
            FlowState.elements.set([]);
            // history.push(`/pipelines/flow/${pipelineId}`);
        } else if (response.r === 'error') {
            closeSnackbar();
            enqueueSnackbar("Can't get flow: " + response.msg, { variant: 'error' });
        } else if (response.errors) {
            response.errors.map((err) => enqueueSnackbar(err.message, { variant: 'error' }));
        } else {
            setElements(response);
            FlowState.elements.set(response);
        }
    };
};

// ------ Custom hooks
export const useGetActiveDeploymentHook = (environmentID, setPipeline) => {
    // GraphQL hook
    const getActiveDeployment = useGetActiveDeployment();

    // URI parameter
    const { deploymentId } = useParams();
    const pipelineID = deploymentId;

    const FlowState = useGlobalFlowState();

    const { enqueueSnackbar, closeSnackbar } = useSnackbar();

    // Get active deployment
    return async () => {
        const response = await getActiveDeployment({ pipelineID, environmentID });

        if (response.r || response.error) {
            closeSnackbar();
            enqueueSnackbar("Can't get active deployment: " + (response.msg || response.r || response.error), { variant: 'error' });
        } else if (response.errors) {
            response.errors.map((err) => enqueueSnackbar(err.message, { variant: 'error' }));
        } else {
            setPipeline(response);
            FlowState.pipelineInfo.set(response);
        }
    };
};

// ----- Utility functions
export function prepareInputForFrontend(input) {
    const edgesInput = [];
    const nodesInput = [];

    if (input && Object.keys(input).length > 0) {
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
                triggerOnline: node.triggerOnline,
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
    }

    return [...edgesInput, ...nodesInput];
}
