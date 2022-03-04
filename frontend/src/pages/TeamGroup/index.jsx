import { Box, Grid, Typography } from '@mui/material';
import { useState, useEffect } from 'react';
import { useSnackbar } from 'notistack';
import { useParams } from 'react-router-dom';
import { useGlobalMeState } from '../../components/Navbar';
import { useGlobalEnvironmentState } from '../../components/EnviromentDropdown';
import { useGetAccessGroup } from '../../graphql/getAccessGroup';
import { useGetEnvironment } from '../../graphql/getEnvironment';
import Details from './Details';
import Control from './Control';
import Permissions from './Permissions';
import Members from './Members';

const TeamGroup = () => {
    // Global user states with hookstate
    const MeData = useGlobalMeState();
    const Environment = useGlobalEnvironmentState();

    // Local state
    const [isGlobalDataLoaded, setIsGlobalDataLoaded] = useState(false);
    const [accessGroup, setAccessGroup] = useState('');
    const [accessGroupEnvironmentName, setAccessGroupEnvironmentName] = useState('');

    // URI parameter
    const { accessId } = useParams();

    // Custom Hook
    const getAccessGroup = useGetAccessGroup_(Environment.id.get(), MeData.user_id.get(), accessId, setAccessGroup, setAccessGroupEnvironmentName);

    // Check if global data is loaded
    useEffect(() => {
        if (Environment.id.get() && MeData.get()) {
            setIsGlobalDataLoaded(true);
            // getAccessGroup();
        }
    }, [Environment, MeData]);

    // Get access group data on load
    useEffect(() => {
        getAccessGroup();

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return (
        <Box className="page" width="83%">
            <Grid container alignItems="flex-start" flexDirection="column">
                <Typography component="h2" variant="h2" color="text.primary">
                    Team {'>'} Access group {'>'} {accessGroup.Name}
                </Typography>

                <Typography variant="subtitle2" mt=".20rem">
                    Environment: {accessGroupEnvironmentName}
                </Typography>
            </Grid>

            {isGlobalDataLoaded && accessGroup ? (
                <Grid container mt="2.56rem" alignItems="flex-start" justifyContent="space-between">
                    <Grid item sx={{ flex: 1 }}>
                        <Details environmentId={Environment.id.get()} accessGroup={accessGroup} getAccessGroup={getAccessGroup} />

                        <Box sx={{ margin: '2.45rem 0', borderTop: 1, borderColor: 'divider' }}></Box>

                        <Control environmentId={Environment.id.get()} accessGroup={accessGroup} getAccessGroup={getAccessGroup} />
                    </Grid>
                    <Grid item sx={{ flex: 2.2, display: 'flex', alignItems: 'center', justifyContent: 'center', flexDirection: 'column' }}>
                        <Permissions environmentId={accessGroup.EnvironmentID} />
                    </Grid>
                    <Grid item sx={{ flex: 1 }}>
                        <Members environmentId={accessGroup.EnvironmentID} />
                    </Grid>
                </Grid>
            ) : null}
        </Box>
    );
};

export default TeamGroup;

// ----------- Custom Hook

const useGetAccessGroup_ = (environmentID, userID, access_group_id, setAccessGroup, setAccessGroupEnvironmentName) => {
    // GraphQL hooks
    const getAccessGroup = useGetAccessGroup();
    const getEnvironment = useGetEnvironment();

    const { enqueueSnackbar, closeSnackbar } = useSnackbar();

    // Get access group data
    return async () => {
        let response = await getAccessGroup({ environmentID, userID, access_group_id });

        if (response.r === 'error') {
            closeSnackbar();
            enqueueSnackbar("Can't get access group data: " + response.msg, {
                variant: 'error',
            });
        } else if (response.errors) {
            response.errors.map((err) => enqueueSnackbar(err.message, { variant: 'error' }));
        } else {
            setAccessGroup(response);

            // Get Access group environment name
            response = await getEnvironment({ environment_id: response.EnvironmentID });

            if (response.r === 'error') {
                closeSnackbar();
                enqueueSnackbar("Can't get access group environment name: " + response.msg, {
                    variant: 'error',
                });
            } else if (response.errors) {
                response.errors.map((err) => enqueueSnackbar(err.message, { variant: 'error' }));
            } else {
                setAccessGroupEnvironmentName(response.name);
            }
        }
    };
};
