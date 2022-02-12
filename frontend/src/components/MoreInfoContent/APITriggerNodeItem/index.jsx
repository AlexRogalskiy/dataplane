import { MenuItem } from '@mui/material';
import { useGlobalFlowState } from '../../../pages/Flow';

const ApiTriggerNodeItem = (props) => {
    const FlowState = useGlobalFlowState();

    const handleDeleteElement = () => {
        FlowState.triggerDelete.set(FlowState.triggerDelete.get() + 1);
        props.handleCloseMenu();
    };

    return (
        <>
            <MenuItem sx={{ color: 'cyan.main' }} onClick={() => props.handleCloseMenu()}>
                API
            </MenuItem>
            <MenuItem sx={{ color: 'error.main' }} onClick={handleDeleteElement}>
                Delete
            </MenuItem>
        </>
    );
};

export default ApiTriggerNodeItem;
