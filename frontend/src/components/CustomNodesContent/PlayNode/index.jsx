import { faPlayCircle } from '@fortawesome/free-regular-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Grid, Typography } from '@mui/material';
import { Box } from '@mui/system';
import { useEffect, useState } from 'react';
import { Handle } from 'react-flow-renderer';
import { useGlobalFlowState } from '../../../pages/Flow';
import customNodeStyle from '../../../utils/customNodeStyle';
import PlayTriggerNodeItem from '../../MoreInfoContent/PlayTriggerNodeItem';
import MoreInfoMenu from '../../MoreInfoMenu';

const PlayNode = () => {
    // Global state
    const FlowState = useGlobalFlowState();

    const [isEditorPage, setIsEditorPage] = useState(false);
    const [isSelected, setIsSelected] = useState(false);

    useEffect(() => {
        setIsEditorPage(FlowState.isEditorPage.get());
        setIsSelected(FlowState.selectedElementId.get() === 'djdsfjdf');
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    useEffect(() => {
        setIsEditorPage(FlowState.isEditorPage.get());
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [FlowState.isEditorPage.get()]);

    useEffect(() => {
        setIsSelected(FlowState.selectedElementId.get() === 'djdsfjdf');
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [FlowState.selectedElementId.get()]);

    return (
        <Box sx={{ ...customNodeStyle, border: `${isSelected ? '3px solid #FF5722' : '3px solid #c4c4c4'}` }}>
            <Handle type="source" position="right" id="play" style={{ backgroundColor: 'red', right: -0.7 }} />
            <Grid container alignItems="flex-start" wrap="nowrap">
                <Box component={FontAwesomeIcon} fontSize={19} color="secondary.main" icon={faPlayCircle} />
                <Grid item ml={1.5} textAlign="left">
                    <Typography fontSize={11} fontWeight={900}>
                        Play trigger
                    </Typography>

                    <Typography fontSize={10} mt={1}>
                        Press play to run
                    </Typography>
                </Grid>

                {isEditorPage && (
                    <Grid position="absolute" bottom={2} right={9} container wrap="nowrap" width="auto" alignItems="center" justifyContent="space-between">
                        <Box mt={2}>
                            <MoreInfoMenu iconHorizontal iconColor="#0073C6" iconColorDark="#0073C6" iconSize={19} noPadding>
                                <PlayTriggerNodeItem />
                            </MoreInfoMenu>
                        </Box>
                    </Grid>
                )}
            </Grid>
        </Box>
    );
};

export default PlayNode;
