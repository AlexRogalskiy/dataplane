import { Box, Container, Grid, Typography } from '@mui/material';
import { useState } from 'react';
import Lottie from 'react-lottie';
import Detail from '../assets/animations/detail.json';
import GetStartedForm from '../components/GetStartedForm';
import SetupLoader from '../components/SetupLoader';

const GetStarted = () => {
    const defaultOptions = {
        loop: true,
        autoplay: true,
        animationData: Detail,
        rendererSettings: {
            preserveAspectRatio: 'xMidYMid slice',
        },
    };

    const [isNext, setIsNext] = useState(false);

    return (
        <Box className="get-started" height="100vh" sx={{ overflowX: 'hidden' }}>
            <Box position="relative" width="100%" sx={{ top: 0, left: 0, right: 0, border: 1, borderColor: 'divider' }}>
                <Typography component="h1" variant="h1" color="secondary" fontWeight={700} style={{ padding: '1.5rem 0 1.5rem .75rem' }}>
                    Dataplane
                </Typography>
            </Box>

            <Container sx={{ mt: 12, mb: 6 }}>
                <Grid container alignItems="flex-start" justifyContent="center">
                    <Grid item flex={1}>
                        <Typography component="h2" fontWeight="700" fontSize="1.93rem" color="text.primary">
                            First time setup
                        </Typography>

                        {isNext ? <SetupLoader /> : <GetStartedForm handleNext={() => setIsNext(true)} />}
                    </Grid>

                    <Grid item flex={2}>
                        <Lottie options={defaultOptions} height={350} width={500} />
                    </Grid>
                </Grid>
            </Container>
        </Box>
    );
};

export default GetStarted;
