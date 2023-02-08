Feature: Process Mini-Ride event's

    As sso
    I want to recive and process mini-ride event's
    So that I will sync new and updated onde driver's

    Scenario: ride mini successfull sync with sso
        Given there are the following user data on sso
            | id                                   | first_name | middle_name | last_name | driverId                             | phone        | profile_picture                                                                                                                       | status    |
            | 06eb340a-862a-4dd0-8a3f-5e4c1f767d3d | abebe      | kebede      | teshome   | aaa5eec3-75d2-4a96-b917-1abda059ec1d | 251944123345 | image                                                                                                                                 | ACTIVE    |
            | 495f6800-dd63-49e2-9809-107076ed2c72 | Surafel    | Zerihun     | Surafel   | 0333b19d-9a8e-4597-95f2-cd2379504c36 | 251967968549 | https://onde-images.s3.amazonaws.com/profile/2021-06-08/0333b19d-9a8e-4597-95f2-cd2379504c36-bfe7b669-40d9-4d03-9a8f-4d78feb93708.png | ACTIVE    |
            | 3088d463-83f6-4a33-94b0-5fcf5b471052 | Genet      | Gezahegn    | Erkita    | 92b51689-3595-4a85-8eeb-5bf3b28a9cbd | 251924301998 | https://onde-images.s3.amazonaws.com/account/2020-08-19/eba18108-23a3-4074-b0fd-ca102f523b2b-78b68c88-8d0b-48c1-9169-1fbf0ea2b8e0.png | ACTIVE    |
            | 19e1b400-3101-49b4-8e04-f57102cb1edb | Bisrat     | Jemal       | Ebrahim   | 4b32a924-e479-4c3d-8568-cfbeabd1ab56 | 251923787979 | https://onde-images.s3.amazonaws.com/account/2020-12-17/fbfc3588-39ce-4f68-9765-2167554c780b-6d9fbafa-b9f5-4409-b7da-bae4435afcc1.png | SUSPENDED |
        And  mini ride streamed the following event's
            | event   | id                                   | full_name               | driver_license | driver_id                            | phone        | profile_picture                                                                                                                       | status | swap_phones                |
            | UPDATE  | 06eb340a-862a-4dd0-8a3f-5e4c1f767d3d | abi kebede teshome      | ab12333        | aaa5eec3-75d2-4a96-b917-1abda059ec1d | 251944123344 | my_image                                                                                                                              | ACTIVE |                            |
            | PROMOTE | 495f6800-dd63-49e2-9809-107076ed2c72 | Surafel Zerihun Surafel | ab12322        | 0333b19d-9a8e-4597-95f2-cd2379504c36 | 251967968549 | https://onde-images.s3.amazonaws.com/profile/2021-06-08/0333b19d-9a8e-4597-95f2-cd2379504c36-bfe7b669-40d9-4d03-9a8f-4d78feb93708.png | ACTIVE |                            |
            | CREATE  | bf576aa8-2945-4e8f-9744-74f1ee5cd7d7 | Yared Amare Sitotaw     | ab12311        | a383e5e1-8d5a-421b-a13c-d3f2b5de4e32 | 251911991471 | https://onde-images.s3.amazonaws.com/account/2020-10-14/fe98ff4e-1239-4ba3-a524-f3ba19d434bf-48e0f257-3f12-4e17-9dcb-34d3b9f3ec1b.png | ACTIVE |                            |
            | UPDATE  | 3088d463-83f6-4a33-94b0-5fcf5b471052 | Genet Gezahegn Erkita   | ab12344        | 92b51689-3595-4a85-8eeb-5bf3b28a9cbd | 251923787979 | https://onde-images.s3.amazonaws.com/account/2020-08-19/eba18108-23a3-4074-b0fd-ca102f523b2b-78b68c88-8d0b-48c1-9169-1fbf0ea2b8e0.png | ACTIVE | 251924301998, 251923787979 |
        When I process those event's
        Then they will have effect on following sso user's
            | id                                   | first_name | middle_name | last_name | phone        | profile_picture                                                                                                                       | status |
            | 06eb340a-862a-4dd0-8a3f-5e4c1f767d3d | abi        | kebede      | teshome   | 251944123344 | my_image                                                                                                                              | ACTIVE |
            | bf576aa8-2945-4e8f-9744-74f1ee5cd7d7 | Yared      | Amare       | Sitotaw   | 251911991471 | https://onde-images.s3.amazonaws.com/account/2020-10-14/fe98ff4e-1239-4ba3-a524-f3ba19d434bf-48e0f257-3f12-4e17-9dcb-34d3b9f3ec1b.png | ACTIVE |
            | 3088d463-83f6-4a33-94b0-5fcf5b471052 | Genet      | Gezahegn    | Erkita    | 251923787979 | https://onde-images.s3.amazonaws.com/account/2020-08-19/eba18108-23a3-4074-b0fd-ca102f523b2b-78b68c88-8d0b-48c1-9169-1fbf0ea2b8e0.png | ACTIVE |

