import { faClock } from '@fortawesome/free-regular-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Grid, Typography, useTheme } from '@mui/material';
import { Box } from '@mui/system';
import { useEffect, useState } from 'react';
import { Handle, useUpdateNodeInternals } from 'react-flow-renderer';
import { useGlobalFlowState } from '../../../pages/Flow';
import customNodeStyle from '../../../utils/customNodeStyle';
import { customSourceConnected, customSourceHandle, customSourceHandleDragging } from '../../../utils/handleStyles';
import ScheduleTriggerNodeItem from '../../MoreInfoContent/ScheduleTriggerNodeItem';
import MoreInfoMenu from '../../MoreInfoMenu';

const ScheduleNode = (props) => {
    const [isRunning, setIsRunning] = useState(false);
    const [isEditorPage, setIsEditorPage] = useState(false);
    const [, setIsSelected] = useState(false);

    // Theme hook
    const theme = useTheme();

    // Update node hook
    // (This will fix the line not having the right size when the node connects to other nodes)
    const updateNodeInternals = useUpdateNodeInternals();

    const FlowState = useGlobalFlowState();

    useEffect(() => {
        setIsRunning(FlowState.isRunning.get());
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [FlowState.isRunning.get()]);

    useEffect(() => {
        setIsEditorPage(FlowState.isEditorPage.get());
        setIsSelected(FlowState.selectedElement.get()?.id === props.id);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    useEffect(() => {
        setIsEditorPage(FlowState.isEditorPage.get());
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [FlowState.isEditorPage.get()]);

    useEffect(() => {
        setIsSelected(FlowState.selectedElement.get()?.id === props.id);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [FlowState.selectedElement.get()]);

    const onConnect = (params) => {
        FlowState.elementsWithConnection.set([...FlowState.elementsWithConnection.get(), params.source]);
        updateNodeInternals(props.id);
    };

    return (
        <Box sx={{ ...customNodeStyle, border: isRunning ? '3px solid #76A853' : '3px solid #c4c4c4' }}>
            <Handle
                onConnect={onConnect}
                type="source"
                position="right"
                id="schedule"
                style={
                    FlowState.isDragging.get()
                        ? customSourceHandleDragging
                        : FlowState.elementsWithConnection.get()?.includes(props.id)
                        ? customSourceConnected(theme.palette.mode)
                        : customSourceHandle(theme.palette.mode)
                }
            />
            <Grid container alignItems="flex-start" wrap="nowrap">
                <Box component={FontAwesomeIcon} fontSize={19} color="secondary.main" icon={faClock} />
                <Grid item ml={1.5} textAlign="left">
                    <Typography fontSize={11} fontWeight={900}>
                        Schedule trigger
                    </Typography>

                    <Typography fontSize={10} mt={1}>
                        Every 5 minutes
                    </Typography>
                </Grid>
            </Grid>

            {isEditorPage && (
                <Grid position="absolute" bottom={2} right={9} container wrap="nowrap" width="auto" alignItems="center" justifyContent="space-between">
                    <Box mt={2}>
                        <MoreInfoMenu iconHorizontal iconColor="#0073C6" iconColorDark="#0073C6" iconSize={19} noPadding>
                            <ScheduleTriggerNodeItem />
                        </MoreInfoMenu>
                    </Box>
                </Grid>
            )}
        </Box>
    );
};

export default ScheduleNode;
