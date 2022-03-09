import { Autocomplete, Grid, TextField } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { useGlobalRunState } from './useWebSocket';
import { useGetPipelineRuns } from '../../graphql/getPipelineRuns';
import { useParams } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import { usePipelineTasksRun } from '../../graphql/getPipelineTasksRun';
import { displayTimer } from './Timer';

export default function RunsDropdown({ environmentID, setElements, setPrevRunTime }) {
    // Global states
    const RunState = useGlobalRunState();

    // Local state
    const [selectedRun, setSelectedRun] = useState();
    const [runs, setRuns] = useState([]);

    // GraphQL hooks
    const getPipelineRuns = useGetPipelineRunsHook(environmentID, setRuns, setSelectedRun);
    const getPipelineTasksRun = usePipelineTasksRunHook();

    // Get pipeline runs on load and environment change and after each run.
    useEffect(() => {
        getPipelineRuns();

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [environmentID, RunState.run_id.get()]);

    // Get pipeline runs on trigger.
    useEffect(() => {
        if (RunState.pipelineRunsTrigger.get() < 2) return;
        getPipelineRuns();

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [RunState.pipelineRunsTrigger.get()]);

    // Update elements on run dropdown change
    useEffect(() => {
        if (!selectedRun) return;
        setElements(selectedRun.run_json);
        getPipelineTasksRun(selectedRun.run_id, environmentID);

        // Set timer on dropdown change. Works only for runs returned from pipeline runs.
        if (selectedRun.ended_at) {
            setPrevRunTime(displayTimer(selectedRun.created_at, selectedRun.ended_at));
        }

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedRun]);

    return (
        <Grid item alignItems="center" display="flex" width={520}>
            {selectedRun || runs.length === 0 ? (
                <Autocomplete
                    id="run_autocomplete"
                    onChange={(event, newValue) => {
                        setSelectedRun(newValue);
                    }}
                    value={selectedRun}
                    disableClearable
                    sx={{ minWidth: '520px' }}
                    options={runs}
                    getOptionLabel={(a) => formatDate(a.created_at) + ' - ' + a.run_id}
                    renderInput={(params) => <TextField {...params} label="Run" id="run" size="small" sx={{ fontSize: '.75rem', display: 'flex' }} />}
                />
            ) : null}
        </Grid>
    );
}

// ------ Custom hook
export const useGetPipelineRunsHook = (environmentID, setRuns, setSelectedRun) => {
    // GraphQL hook
    const getPipelineRuns = useGetPipelineRuns();

    const RunState = useGlobalRunState();

    // URI parameter
    const { pipelineId } = useParams();

    const { enqueueSnackbar, closeSnackbar } = useSnackbar();

    // Get members
    return async () => {
        const response = await getPipelineRuns({ pipelineID: pipelineId, environmentID });

        if (response.length === 0) {
            setRuns([]);
        } else if (response.r === 'error') {
            closeSnackbar();
            enqueueSnackbar("Can't get flow: " + response.msg, { variant: 'error' });
        } else if (response.errors) {
            response.errors.map((err) => enqueueSnackbar(err.message, { variant: 'error' }));
        } else {
            setRuns(response);
            setSelectedRun(response[0]);
        }
    };
};

export const usePipelineTasksRunHook = () => {
    // GraphQL hook
    const getPipelineTasksRun = usePipelineTasksRun();

    // URI parameter
    const { pipelineId } = useParams();

    const RunState = useGlobalRunState();

    const { enqueueSnackbar, closeSnackbar } = useSnackbar();

    // Update pipeline flow
    return async (runID, environmentID) => {
        if (!runID) return;

        const response = await getPipelineTasksRun({ pipelineID: pipelineId, runID, environmentID });

        if (response.r === 'Unauthorized') {
            closeSnackbar();
            enqueueSnackbar(`Can't update flow: ${response.r}`, { variant: 'error' });
        } else if (response.errors) {
            response.errors.map((err) => enqueueSnackbar(err.message, { variant: 'error' }));
        } else {
            // Keeping only start_id and run_id and removing rest of the nodes before adding this run's nodes.
            const keep = { start_id: RunState.start_dt.get(), run_id: RunState.run_id.get(), pipelineRunsTrigger: RunState.pipelineRunsTrigger.get() };
            response.map((a) => (keep[a.node_id] = { status: a.status, end_dt: a.end_dt, start_dt: a.start_dt }));
            RunState.set(keep);
        }
    };
};

// ----- Utility function
function formatDate(date) {
    if (!date) return;
    date = new Date(date);
    let day = new Intl.DateTimeFormat('en', { day: 'numeric' }).format(date);
    let monthYear = new Intl.DateTimeFormat('en', { year: 'numeric', month: 'short' }).format(date);
    let time = new Intl.DateTimeFormat('en', { hourCycle: 'h23', hour: '2-digit', minute: 'numeric', second: 'numeric' }).format(date);
    return `${day} ${monthYear} ${time}`;
}
