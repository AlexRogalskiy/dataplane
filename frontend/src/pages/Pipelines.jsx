import { Box, Grid, Typography, Chip, Avatar, TextField, InputAdornment, MenuItem } from '@mui/material';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEllipsisV, faChartBar } from '@fortawesome/free-solid-svg-icons'
import Search from '../components/Search';
import CustomChip from '../components/CustomChip';
import PipelineTable from '../components/TableContent/PipelineTable';

const Pipelines = () => {
    return (
        <Box className="page">
            <Grid container alignItems="center" justifyContent="space-between">
                <Typography component="h2" variant="h2" color="text.primary">
                    Pipelines
                </Typography>
                <FontAwesomeIcon icon={faEllipsisV} />
            </Grid>

            <Grid container mt={4} direction="row" alignItems="center" justifyContent="flex-start" sx={{ width: { xl: "85%" } }}>
                <Grid item display="flex" alignItems="center" sx={{ alignSelf: "center" }}>
                    <CustomChip amount={2} label="Pipelines" margin={1} customColor="orange" />
                    <CustomChip amount={2} label="Succeeded" margin={1} customColor="green" />
                    <CustomChip amount={2} label="Failed" margin={1} customColor="red" />
                    <CustomChip amount={2} label="Workers online" margin={2} customColor="purple" />

                </Grid>

                <Grid item display="flex" alignItems="center" sx={{ alignSelf: "center", flex: 1 }}>
                    <Box flex={1.2} width="100%" flexGrow={1.2}>
                        <Search placeholder="Find a pipeline" />
                    </Box>
                    <TextField 
                        label="Last 48 hours"
                        id="last"
                        select
                        size="small"
                        required
                        sx={{ ml: 2, flex: 1 }}
                    >
                        <MenuItem value="24">
                            Last 24 hours
                        </MenuItem>
                    </TextField>
                </Grid>

                <PipelineTable />
            </Grid>
        </Box>
    )
}

export default Pipelines;